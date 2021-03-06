// The MIT License (MIT)
//
// Copyright © 2018 Sven Agneessens <sven.agneessens@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package input

import (
	"fmt"
	"sync"

	"github.com/apex/log"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type MQTT struct {
	options  MQTTOptions
	client   paho.Client
	cOptions *paho.ClientOptions

	subscriptions map[string]bool

	Incoming chan paho.Message
	Done     chan struct{}

	sync.Mutex
}

type MQTTOptions struct {
	Server   string
	Username string
	Password string
	QoS      int
	ClientID string
	Debug    bool
}

type debugLogger struct{}

func New(options MQTTOptions) *MQTT {
	if options.Debug {
		paho.DEBUG = debugLogger{}
	}

	return &MQTT{
		options: options,
	}
}

func (m *MQTT) Connect() error {
	m.Lock()
	defer m.Unlock()

	var err error

	if m.options.QoS < 0 || m.options.QoS > 2 {
		return errors.Errorf("[MQTT] invalid QoS: %v", m.options.QoS)
	}

	m.cOptions, err = m.createOptions()
	if err != nil {
		return errors.Wrap(err, "[MQTT] error creating options")
	}

	var reconnecting bool

	m.cOptions.SetConnectionLostHandler(func(client paho.Client, err error) {
		log.Warnf("[MQTT] disconnected (%s), reconnecting...", err)
		reconnecting = true
	})

	m.cOptions.SetOnConnectHandler(func(client paho.Client) {
		log.Info("[MQTT] connected")
		if reconnecting {
			for topic, on := range m.subscriptions {
				if on {
					log.Debugf("[MQTT] re-subscribing to topic: %s", topic)
					m.Subscribe(topic)
				}
			}
			reconnecting = false
		}
	})

	m.client = paho.NewClient(m.cOptions)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "[MQTT] error connecting")
	}

	m.subscriptions = make(map[string]bool)
	m.Incoming = make(chan paho.Message)
	m.Done = make(chan struct{})

	return nil
}

func (m *MQTT) createOptions() (*paho.ClientOptions, error) {
	cOptions := paho.NewClientOptions()

	cOptions.AddBroker(m.options.Server)

	if m.options.Username != "" && m.options.Password != "" {
		cOptions.SetUsername(m.options.Username)
		cOptions.SetPassword(m.options.Password)
	}

	cOptions.SetClientID(m.options.ClientID)

	return cOptions, nil
}

func (m *MQTT) Subscribe(topic string) error {
	if m.client != nil && m.client.IsConnected() {
		if token := m.client.Subscribe(topic, byte(m.options.QoS), m.onReceive); token.Wait() && token.Error() != nil {
			return errors.Wrapf(token.Error(), "[MQTT] error subscribing to %s", topic)
		}

		m.subscriptions[topic] = true

		log.Infof("[MQTT] subscribing to topic: %s", topic)

		return nil
	} else {
		return errors.New("[MQTT] trying to subscribe while not connected")
	}
}

func (m *MQTT) onReceive(_ paho.Client, message paho.Message) {
	m.Incoming <- message
}

func (m *MQTT) Unsubscribe(topics ...string) error {
	if m.client != nil && m.client.IsConnected() {
		if token := m.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
			return errors.Wrapf(token.Error(), "[MQTT] error unsubscribing from: %s", topics)
		}

		for _, topic := range topics {
			m.subscriptions[topic] = false

			log.Infof("[MQTT] un-subscribing from topic: %s", topic)
		}

		return nil
	} else {
		return errors.New("[MQTT] trying to unsubscribe while not connected")
	}
}

func (m *MQTT) Close() {
	m.Lock()
	defer m.Unlock()

	if m.client != nil && m.client.IsConnected() {
		close(m.Done)
		m.client.Disconnect(250)
		log.Info("[MQTT] disconnected")
	}
}

func (l debugLogger) Println(v ...interface{}) {
	log.Debugf("[MQTT Debug] %s", fmt.Sprintln(v...))
}

func (l debugLogger) Printf(format string, v ...interface{}) {
	log.Debugf("[MQTT Debug] %s", fmt.Sprintf(format, v...))
}

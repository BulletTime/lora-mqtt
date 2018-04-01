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
	"sync"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type MQTT struct {
	options  MQTTOptions
	client   paho.Client
	cOptions *paho.ClientOptions

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
}

func New(options MQTTOptions) *MQTT {
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

	m.client = paho.NewClient(m.cOptions)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "[MQTT] error connecting")
	}

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
	}
}

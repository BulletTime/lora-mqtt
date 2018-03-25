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
	"testing"
	"fmt"
	"github.com/bullettime/lora-mqtt/util"
)

var options = MQTTOptions{
	Server: "tcp://localhost:1883",
	Username: "",
	Password: "",
	QoS: 0,
	ClientID: fmt.Sprintf("lora-mqtt-%s", util.RandomString(4)),
}

func TestNew(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test for default values")
	}

	mqtt := New(options)

	if mqtt.options.Server != "tcp://localhost:1883" {
		t.Error("invalid default server")
	}

	if mqtt.options.Username != "" {
		t.Error("invalid default username")
	}

	if mqtt.options.Password != "" {
		t.Error("invalid default Password")
	}

	if mqtt.options.QoS != 0 {
		t.Error("invalid default qos")
	}
}

func TestMQTT_Connect(t *testing.T) {
	mqtt := New(options)

	if err := mqtt.Connect(); err != nil {
		t.Error(err)
	}

	mqtt.Close()
}

func TestMQTT_Connect2(t *testing.T) {
	mqtt := New(options)

	mqtt.options.QoS = 3

	if err := mqtt.Connect(); err == nil {
		t.Error("invalid qos at connection")
	}

	mqtt.options.QoS = -1

	if err := mqtt.Connect(); err == nil {
		t.Error("invalid qos at connection")
	}

	mqtt.Close()
}

func TestMQTT_Subscribe(t *testing.T) {
	mqtt := New(options)

	if err := mqtt.Connect(); err != nil {
		t.Error(err)
	}

	if err := mqtt.Subscribe("/test"); err != nil {
		t.Error(err)
	}

	mqtt.Close()
}

func TestMQTT_Subscribe2(t *testing.T) {
	mqtt := New(options)

	if err := mqtt.Subscribe("/test"); err == nil {
		t.Error("subscribing while not connected")
	}
}

func TestMQTT_Unsubscribe(t *testing.T) {
	mqtt := New(options)

	if err := mqtt.Connect(); err != nil {
		t.Error(err)
	}

	if err := mqtt.Subscribe("/test"); err != nil {
		t.Error(err)
	}

	if err := mqtt.Unsubscribe("/test"); err != nil {
		t.Error(err)
	}

	mqtt.Close()
}

func TestMQTT_Unsubscribe2(t *testing.T) {
	mqtt := New(options)

	if err := mqtt.Connect(); err != nil {
		t.Error(err)
	}

	if err := mqtt.Unsubscribe("/test"); err != nil {
		t.Error(err)
	}

	mqtt.Close()
}

func TestMQTT_Unsubscribe3(t *testing.T) {
	mqtt := New(options)

	if err := mqtt.Unsubscribe("/test"); err == nil {
		t.Error("unsubscribing while not connected")
	}
}

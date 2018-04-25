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

package dingnetjson

import (
	"testing"
)

const (
	name = "test"
)

func TestNew(t *testing.T) {
	p, err := New(name)
	if err != nil {
		t.Error(err)
	}
	if p.(*dingnetParser).MetricName != name {
		t.Error("metric name should be initialized initialized")
	}

	p, err = New("")
	if err == nil {
		t.Error("empty metric name should give an error")
	}
}

//func TestTtnParser_Parse(t *testing.T) {
//	p, err := New(parser.LocationData)
//	if err != nil {
//		t.Error(err)
//	}
//
//	metrics, err := p.Parse([]byte(jsonMessage))
//	if err != nil {
//		t.Error(err)
//	}
//
//	if len(metrics) != 1 {
//		t.Error("should only have 1 metric")
//	}
//
//	metric := metrics[0]
//
//	if !(metric.HasTag("device_id") && metric.HasTag("frequency") && metric.HasTag("data_rate") &&
//		metric.HasTag("power") && metric.HasTag("latitude") && metric.HasTag("longitude") &&
//		metric.HasTag("gateway_id")) {
//		t.Error("missing one or more tags")
//	}
//
//	if !(metric.HasField("size") && metric.HasField("rssi") && metric.HasField("snr")) {
//		t.Error("missing one or more fields")
//	}
//}
//
//func TestTtnParser_Parse2(t *testing.T) {
//	p, err := New(name)
//	if err != nil {
//		t.Error(err)
//	}
//
//	metrics, err := p.Parse([]byte(jsonMessage))
//	if err != nil {
//		t.Error(err)
//	}
//
//	if len(metrics) != 1 {
//		t.Error("should only have 1 metric")
//	}
//
//	metric := metrics[0]
//
//	if !(metric.HasTag("device_id") && metric.HasTag("frequency") && metric.HasTag("data_rate") &&
//		metric.HasTag("rssi") && metric.HasTag("snr") && metric.HasTag("gateway_id")) {
//		t.Error("missing one or more tags")
//	}
//
//	if !(metric.HasField("size") && metric.HasField("lat") && metric.HasField("lon") &&
//		metric.HasField("pwr")) {
//		t.Error("missing one or more fields")
//	}
//}
//
//func TestTtnParser_Parse3(t *testing.T) {
//	p, err := New(name)
//	if err != nil {
//		t.Error(err)
//	}
//
//	_, err = p.Parse([]byte(jsonMessageNoGateways))
//	if err == nil {
//		t.Error("should not be able to parse a json message without gateways info")
//	}
//}

func TestTtnParser_SetDefaultTags(t *testing.T) {
	p, err := New(name)
	if err != nil {
		t.Error(err)
	}

	tags := map[string]string{
		"test": "a",
	}

	p.SetDefaultTags(tags)

	if v, ok := p.(*dingnetParser).DefaultTags["test"]; !ok {
		t.Error("default tags is missing key 'test'")
	} else {
		if v != "a" {
			t.Error("default tags has wrong value for key 'test'")
		}
	}
}

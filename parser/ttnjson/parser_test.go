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

package ttnjson

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"testing"
)

const (
	name        = "test"
	jsonMessage = `{
  "app_id": "lora_coverage_mapping",
  "dev_id": "sodaq_one_gps_1",
  "hardware_serial": "003017737253C1D7",
  "port": 1,
  "counter": 7,
  "payload_raw": "B8hBALggAQ==",
  "payload_fields": {
    "lat": 51.0017,
	"lon": 4.7136,
	"pwr": 1
  },
  "metadata": {
	"airtime": 1318912000,
	"time": "2018-03-13T19:21:22.827671626Z",
	"frequency": 868.3,
	"modulation": "LORA",
	"data_rate": "SF12BW125",
	"coding_rate": "4/5",
	"gateways": [
	  {
		"gateway_id": "eui-008000000000b88d",
		"timestamp": 3239248428,
		"time": "2018-03-13T19:19:30.671066Z",
		"channel": 1,
		"rssi": -84,
		"snr": 8,
		"rf_chain": 1
	  }
	]
  }
}`
	jsonMessageNoGateways = `{
  "app_id": "lora_coverage_mapping",
  "dev_id": "sodaq_one_gps_1",
  "hardware_serial": "003017737253C1D7",
  "port": 1,
  "counter": 7,
  "payload_raw": "B8hBALggAQ==",
  "payload_fields": {
    "lat": 51.0017,
	"lon": 4.7136,
	"pwr": 1
  },
  "metadata": {
	"airtime": 1318912000,
	"time": "2018-03-13T19:21:22.827671626Z",
	"frequency": 868.3,
	"modulation": "LORA",
	"data_rate": "SF12BW125",
	"coding_rate": "4/5",
  }
}`
)

func TestNew(t *testing.T) {
	p, err := New(name)
	if err != nil {
		t.Error(err)
	}
	if p.(*ttnParser).MetricName != name {
		t.Error("metric name should be initialized initialized")
	}

	p, err = New("")
	if err == nil {
		t.Error("empty metric name should give an error")
	}
}

func TestPayload_MarshalJSON(t *testing.T) {
	p := payload{
		size:  4,
		bytes: []byte{0x0, 0x1, 0x2, 0x3},
	}

	out, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if out[0] != '"' || out[len(out)-1] != '"' {
		t.Fatal("payload should be quoted")
	}

	decoded, err := base64.StdEncoding.DecodeString(string(out[1 : len(out)-1]))
	if err != nil {
		t.Fatal(err)
	}

	if len(decoded) != p.size {
		t.Error("incorrect size after marshalling")
	}

	if bytes.Compare(p.bytes, decoded) != 0 {
		t.Error("bytes not matching before and after marshalling")
	}
}

func TestPayload_MarshalJSON2(t *testing.T) {
	p := payload{}

	out, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if out[0] != '"' || out[len(out)-1] != '"' {
		t.Fatal("payload should be quoted")
	}

	decoded, err := base64.StdEncoding.DecodeString(string(out[1 : len(out)-1]))
	if err != nil {
		t.Fatal(err)
	}

	if len(decoded) != p.size {
		t.Error("incorrect size after marshalling")
	}

	if bytes.Compare(p.bytes, decoded) != 0 {
		t.Error("bytes before and after marshalling not matching")
	}
}

func TestPayload_UnmarshalJSON(t *testing.T) {
	p := payload{}
	testBytes := []byte{0x0, 0x1, 0x2, 0x3}

	encoded := base64.StdEncoding.EncodeToString(testBytes)
	encoded = strconv.Quote(encoded)

	err := p.UnmarshalJSON([]byte(encoded))
	if err != nil {
		t.Error(err)
	}

	if len(testBytes) != p.size {
		t.Error("incorrect size after unmarshalling")
	}

	if bytes.Compare(testBytes, p.bytes) != 0 {
		t.Error("bytes before and after marshalling not matching")
	}
}

func TestPayload_UnmarshalJSON2(t *testing.T) {
	p := payload{}
	testBytes := make([]byte, 0)

	encoded := base64.StdEncoding.EncodeToString(testBytes)
	encoded = strconv.Quote(encoded)

	err := p.UnmarshalJSON([]byte(encoded))
	if err != nil {
		t.Error(err)
	}

	if len(testBytes) != p.size {
		t.Error("incorrect size after unmarshalling")
	}

	if bytes.Compare(testBytes, p.bytes) != 0 {
		t.Error("bytes before and after marshalling not matching")
	}
}

func TestTtnParser_Parse(t *testing.T) {
	p, err := New(LocationData)
	if err != nil {
		t.Error(err)
	}

	metrics, err := p.Parse([]byte(jsonMessage))
	if err != nil {
		t.Error(err)
	}

	if len(metrics) != 1 {
		t.Error("should only have 1 metric")
	}

	metric := metrics[0]

	if !(metric.HasTag("device_id") && metric.HasTag("frequency") && metric.HasTag("data_rate") &&
		metric.HasTag("power") && metric.HasTag("latitude") && metric.HasTag("longitude") &&
		metric.HasTag("gateway_id")) {
		t.Error("missing one or more tags")
	}

	if !(metric.HasField("size") && metric.HasField("rssi") && metric.HasField("snr")) {
		t.Error("missing one or more fields")
	}
}

func TestTtnParser_Parse2(t *testing.T) {
	p, err := New(name)
	if err != nil {
		t.Error(err)
	}

	metrics, err := p.Parse([]byte(jsonMessage))
	if err != nil {
		t.Error(err)
	}

	if len(metrics) != 1 {
		t.Error("should only have 1 metric")
	}

	metric := metrics[0]

	if !(metric.HasTag("device_id") && metric.HasTag("frequency") && metric.HasTag("data_rate") &&
		metric.HasTag("rssi") && metric.HasTag("snr") && metric.HasTag("gateway_id")) {
		t.Error("missing one or more tags")
	}

	if !(metric.HasField("size") && metric.HasField("lat") && metric.HasField("lon") &&
		metric.HasField("pwr")) {
		t.Error("missing one or more fields")
	}
}

func TestTtnParser_Parse3(t *testing.T) {
	p, err := New(name)
	if err != nil {
		t.Error(err)
	}

	_, err = p.Parse([]byte(jsonMessageNoGateways))
	if err == nil {
		t.Error("should not be able to parse a json message without gateways info")
	}
}

func TestTtnParser_SetDefaultTags(t *testing.T) {
	p, err := New(name)
	if err != nil {
		t.Error(err)
	}

	tags := map[string]string {
		"test": "a",
	}

	p.SetDefaultTags(tags)

	if v, ok := p.(*ttnParser).DefaultTags["test"]; !ok {
		t.Error("default tags is missing key 'test'")
	} else {
		if v != "a" {
			t.Error("default tags has wrong value for key 'test'")
		}
	}
}

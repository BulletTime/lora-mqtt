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

package influxdb

import (
	"testing"
	"time"

	"github.com/bullettime/lora-mqtt/model"
)

var server = "http://localhost:8086"

var options = InfluxOptions{
	Server:   server,
	Username: "demo",
	Password: "demo",
	Database: "demo",
}

func TestInfluxdb_Connect(t *testing.T) {
	influxdb := New(options)

	if err := influxdb.Connect(); err != nil {
		t.Error(err)
	}

	influxdb.Close()
}

func TestInfluxdb_Connect2(t *testing.T) {
	options2 := options
	options2.Server = "http://localhost:8080"
	influxdb := New(options2)

	if err := influxdb.Connect(); err == nil {
		t.Error("Wrong server should give error")
	}

	influxdb.Close()
}

func TestInfluxdb_Write(t *testing.T) {
	name := "mysensor1"
	tags := make(map[string]string)
	fields := make(map[string]interface{})

	influxdb := New(options)

	tags["testtag"] = "testing"
	fields["value"] = 0.65

	metric, err := model.NewMetric(name, tags, fields, time.Now())
	if err != nil {
		t.Error(err)
	}

	if err := influxdb.Connect(); err != nil {
		t.Error(err)
	}

	if err := influxdb.Write([]model.Metric{metric}); err != nil {
		t.Error(err)
	}

	time.Sleep(500 * time.Millisecond)
	metric2, err := model.NewMetric(name, tags, fields, time.Time{})
	if err != nil {
		t.Error(err)
	}

	if err := influxdb.Write([]model.Metric{metric2}); err != nil {
		t.Error(err)
	}

	influxdb.Close()
}

func TestInfluxdb_Write2(t *testing.T) {
	name := "mysensor2"
	tags := make(map[string]string)
	fields := make(map[string]interface{})

	influxdb := New(options)

	tags["testtag"] = "testing"
	fields["value"] = 0.75

	metric, err := model.NewMetric(name, tags, fields, time.Now())
	if err != nil {
		t.Error(err)
	}

	time.Sleep(100 * time.Millisecond)
	fields["value"] = 0.85

	metric2, err := model.NewMetric(name, tags, fields, time.Now())
	if err != nil {
		t.Error(err)
	}

	if err := influxdb.Connect(); err != nil {
		t.Error(err)
	}

	if err := influxdb.Write([]model.Metric{metric, metric2}); err != nil {
		t.Error(err)
	}

	influxdb.Close()
}

func TestInfluxdb_Close(t *testing.T) {
	influxdb := New(options)

	if err := influxdb.Connect(); err != nil {
		t.Error(err)
	}

	if err := influxdb.Close(); err != nil {
		t.Error(err)
	}
}

func TestInfluxdb_Close2(t *testing.T) {
	influxdb := New(options)

	if err := influxdb.Close(); err != nil {
		t.Error(err)
	}
}

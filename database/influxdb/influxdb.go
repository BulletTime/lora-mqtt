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
	"time"

	"github.com/bullettime/lora-mqtt/database"
	"github.com/bullettime/lora-mqtt/model"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

type influxdb struct {
	client  client.Client
	options InfluxOptions
}

type InfluxOptions struct {
	Server    string
	Username  string
	Password  string
	Database  string
	Precision string
}

func New(options InfluxOptions) database.Database {
	return &influxdb{
		options: options,
	}
}

func (i *influxdb) Connect() error {
	var err error

	config := client.HTTPConfig{
		Addr:     i.options.Server,
		Username: i.options.Username,
		Password: i.options.Password,
	}

	i.client, err = client.NewHTTPClient(config)
	if err != nil {
		return errors.Wrap(err, "[Influxdb] error creating http client")
	}

	_, _, err = i.client.Ping(3 * time.Second)
	if err != nil {
		return errors.Wrap(err, "[Influxdb] error establishing connection")
	}

	return nil
}

func (i *influxdb) Write(metrics []model.Metric) error {
	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.options.Database,
		Precision: i.options.Precision,
	})
	if err != nil {
		return errors.Wrap(err, "[Influxdb] error creating new batch points")
	}

	for _, metric := range metrics {
		point, err := client.NewPoint(metric.Name(), metric.Tags(), metric.Fields(), metric.Time())
		if err != nil {
			continue
		}

		batchPoints.AddPoint(point)
	}
	if err != nil {
		return errors.Wrap(err, "[Influxdb] error creating new point(s)")
	}

	if err := i.client.Write(batchPoints); err != nil {
		return errors.Wrap(err, "[Influxdb] error writing batch points")
	}

	return nil
}

func (i *influxdb) Close() error {
	if i.client != nil {
		return i.client.Close()
	}

	return nil
}

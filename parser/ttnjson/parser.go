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
	"encoding/json"
	"strconv"
	"time"

	"github.com/bullettime/lora-mqtt/model"
	"github.com/bullettime/lora-mqtt/parser"
	"github.com/pkg/errors"
)

type ttnParser struct {
	MetricName  string
	DefaultTags map[string]string
}

type ttnJson struct {
	AppID          string                 `json:"app_id"`
	DevID          string                 `json:"dev_id"`
	HardwareSerial string                 `json:"hardware_serial"`
	Port           int                    `json:"port"`
	Counter        int                    `json:"counter"`
	IsRetry        bool                   `json:"is_retry,omitempty"`
	Confirmed      bool                   `json:"confirmed,omitempty"`
	PayloadRaw     parser.Payload         `json:"payload_raw"`
	PayloadFields  map[string]interface{} `json:"payload_fields,omitempty"`
	Metadata       metadata               `json:"metadata"`
}

type metadata struct {
	Airtime    time.Duration `json:"airtime"`
	Time       time.Time     `json:"time"`
	Frequency  float64       `json:"frequency"`
	Modulation string        `json:"modulation"`
	DataRate   string        `json:"data_rate"`
	BitRate    int           `json:"bit_rate,omitempty"`
	CodingRate string        `json:"coding_rate"`
	Gateways   []gateway     `json:"gateways"`
	Latitude   float64       `json:"latitude,omitempty"`
	Longitude  float64       `json:"longitude,omitempty"`
	Altitude   float64       `json:"altitude,omitempty"`
}

type gateway struct {
	GatewayID string    `json:"gtw_id"`
	Timestamp uint64    `json:"timestamp"`
	Time      time.Time `json:"time"`
	Channel   int       `json:"channel"`
	RSSI      int       `json:"rssi"`
	SNR       float64   `json:"snr"`
	RfChain   int       `json:"rf_chain"`
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Altitude  float64   `json:"altitude,omitempty"`
}

func New(name string) (parser.Parser, error) {
	if len(name) == 0 {
		return nil, errors.New("[TTNParser] name cannot be empty")
	}

	p := ttnParser{
		MetricName: name,
	}

	return &p, nil
}

func (p *ttnParser) Parse(buf []byte) ([]model.Metric, error) {
	var metrics []model.Metric
	var message ttnJson

	err := json.Unmarshal(buf, &message)
	if err != nil {
		return nil, errors.Wrapf(err, "[TTNParser] error unmarshalling byte buffer: %s", string(buf))
	}

	if len(message.Metadata.Gateways) == 0 {
		return nil, errors.New("[TTNParser] wrong number of gateways (0)")
	}

	tags := make(map[string]string, len(p.DefaultTags))
	for k, v := range p.DefaultTags {
		tags[k] = v
	}

	tags["device_id"] = message.DevID
	tags["frequency"] = strconv.FormatFloat(message.Metadata.Frequency, 'f', -1, 64)
	tags["data_rate"] = message.Metadata.DataRate
	if p.MetricName == parser.LocationData {
		power, err := message.PayloadRaw.GetPower()
		if err == nil {
			tags["power"] = strconv.Itoa(int(power))
		}
		lat, lon, err := message.PayloadRaw.GetLocation()
		if err == nil {
			tags["latitude"] = strconv.FormatFloat(lat, 'f', 4, 64)
			tags["longitude"] = strconv.FormatFloat(lon, 'f', 4, 64)
		}
	}

	fields := map[string]interface{}{
		"size": message.PayloadRaw.Size,
	}

	for _, g := range message.Metadata.Gateways {
		metric, err := model.NewMetric(p.MetricName, tags, fields, g.Time)
		if err != nil {
			return nil, errors.Wrap(err, "[TTNParser] error creating metric")
		}

		if p.MetricName == parser.LocationData {
			metric.AddField("rssi", g.RSSI)
			metric.AddField("snr", g.SNR)
		} else {
			metric.AddTag("rssi", strconv.Itoa(g.RSSI))
			metric.AddTag("snr", strconv.FormatFloat(g.SNR, 'f', -1, 64))
			for k, v := range message.PayloadFields {
				metric.AddField(k, v)
			}
		}

		metric.AddTag("gateway_id", g.GatewayID)

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (p *ttnParser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}

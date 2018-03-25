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

package model

import (
	"time"
	"github.com/pkg/errors"
)

type metric struct {
	name string
	tags map[string]string
	fields map[string]interface{}
	time time.Time
}

type Metric interface {
	// Getting data structure functions
	Name() string
	Tags() map[string]string
	Fields() map[string]interface{}
	Time() time.Time

	// Tag functions
	HasTag(key string) bool
	AddTag(key, value string)
	RemoveTag(key string)

	// Field functions
	HasField(key string) bool
	AddField(key string, value interface{})
	RemoveField(key string) error
}

func NewMetric(name string, tags map[string]string, fields map[string]interface{}, time time.Time) (Metric, error) {
	if len(name) == 0 {
		return nil, errors.New("[Metric] missing measurement name")
	}

	if len(fields) == 0 {
		return nil, errors.Errorf("[Metric] %s: missing field(s) (at least one required)")
	}

	m := &metric{
		name: name,
		tags: tags,
		fields: fields,
		time: time,
	}

	return m, nil
}

func (m *metric) Name() string {
	return m.name
}

func (m *metric) Tags() map[string]string {
	return m.tags
}

func (m *metric) Fields() map[string]interface{} {
	return m.fields
}

func (m *metric) Time() time.Time {
	return m.time
}

func (m *metric) HasTag(key string) bool {
	_, ok := m.tags[key]
	return ok
}

func (m *metric) AddTag(key, value string) {
	m.tags[key] = value
}

func (m *metric) RemoveTag(key string) {
	delete(m.tags, key)
}

func (m *metric) HasField(key string) bool {
	_, ok := m.fields[key]
	return ok
}

func (m *metric) AddField(key string, value interface{}) {
	m.fields[key] = value
}

func (m *metric) RemoveField(key string) error {
	if len(m.fields) == 1 {
		return errors.New("[Metric] can't delete last field (at least one required)")
	}
	delete(m.fields, key)
	return nil
}

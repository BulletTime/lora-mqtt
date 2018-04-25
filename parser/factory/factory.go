//go:generate stringer -type=TypeParser
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

package factory

import (
	"github.com/bullettime/lora-mqtt/parser"
	"github.com/bullettime/lora-mqtt/parser/dingnetjson"
	"github.com/bullettime/lora-mqtt/parser/ttnjson"
	"github.com/pkg/errors"
)

type TypeParser int

const (
	TTN TypeParser = iota
	DingNet
)

func GetTypesList() []string {
	var typesList []string
	for i := 0; i < len(_TypeParser_index)-1; i++ {
		typesList = append(typesList, TypeParser(i).String())
	}
	return typesList
}

func CreateParser(typeParser TypeParser, metricName string) (parser.Parser, error) {
	switch typeParser {
	case TTN:
		return ttnjson.New(metricName)
	case DingNet:
		return dingnetjson.New(metricName)
	default:
		return nil, errors.New("[Parser Factory] incorrect parser type")
	}
}

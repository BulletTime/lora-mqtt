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

package parser

import (
	"encoding/base64"
	"strconv"
	"github.com/pkg/errors"
)

var InvalidPayloadError = errors.New("invalid payload")

type Payload struct {
	Size  int
	Bytes []byte
}

func (p Payload) MarshalJSON() ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.EncodedLen(p.Size))
	base64.StdEncoding.Encode(buf, p.Bytes)
	bufStr := strconv.Quote(string(buf))
	return []byte(bufStr), nil
}

func (p *Payload) UnmarshalJSON(data []byte) error {
	dataStr, err := strconv.Unquote(string(data))
	if err != nil {
		return errors.Wrap(err, "[Payload] error unquoting raw Payload (base64)")
	}
	dataBytes := []byte(dataStr)
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(dataBytes)))
	p.Size, err = base64.StdEncoding.Decode(buf, dataBytes)
	if err != nil {
		return errors.Wrap(err, "[Payload] error unmarshalling raw Payload (base64)")
	}
	bufBytes := make([]byte, p.Size)
	copy(bufBytes, buf[:p.Size])
	p.Bytes = bufBytes
	return nil
}

func (p Payload) GetLocation() (float64, float64, error) {
	if !p.IsValidPayload() {
		return 0, 0, InvalidPayloadError
	}

	multiplier := float64(10000)

	latitude := float64(uint32(p.Bytes[2])|uint32(p.Bytes[1])<<8|uint32(p.Bytes[0])<<16) / multiplier
	longitude := float64(uint32(p.Bytes[5])|uint32(p.Bytes[4])<<8|uint32(p.Bytes[3])<<16) / multiplier

	return latitude, longitude, nil
}

func (p Payload) GetPower() (int8, error) {
	if !p.IsValidPayload() || len(p.Bytes) < 7 {
		return 127, InvalidPayloadError
	}

	power := int8(p.Bytes[6])

	return power, nil
}

func (p Payload) IsValidPayload() bool {
	if len(p.Bytes) < 6 || len(p.Bytes) > 7 {
		return false
	} else {
		return true
	}
}

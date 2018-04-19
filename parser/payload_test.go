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
	"testing"
	"encoding/base64"
	"bytes"
	"strconv"
)

func TestPayload_MarshalJSON(t *testing.T) {
	p := Payload{
		Size:  4,
		Bytes: []byte{0x0, 0x1, 0x2, 0x3},
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

	if len(decoded) != p.Size {
		t.Error("incorrect size after marshalling")
	}

	if bytes.Compare(p.Bytes, decoded) != 0 {
		t.Error("bytes not matching before and after marshalling")
	}
}

func TestPayload_MarshalJSON2(t *testing.T) {
	p := Payload{}

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

	if len(decoded) != p.Size {
		t.Error("incorrect size after marshalling")
	}

	if bytes.Compare(p.Bytes, decoded) != 0 {
		t.Error("bytes before and after marshalling not matching")
	}
}

func TestPayload_UnmarshalJSON(t *testing.T) {
	p := Payload{}
	testBytes := []byte{0x0, 0x1, 0x2, 0x3}

	encoded := base64.StdEncoding.EncodeToString(testBytes)
	encoded = strconv.Quote(encoded)

	err := p.UnmarshalJSON([]byte(encoded))
	if err != nil {
		t.Error(err)
	}

	if len(testBytes) != p.Size {
		t.Error("incorrect size after unmarshalling")
	}

	if bytes.Compare(testBytes, p.Bytes) != 0 {
		t.Error("bytes before and after marshalling not matching")
	}
}

func TestPayload_UnmarshalJSON2(t *testing.T) {
	p := Payload{}
	testBytes := make([]byte, 0)

	encoded := base64.StdEncoding.EncodeToString(testBytes)
	encoded = strconv.Quote(encoded)

	err := p.UnmarshalJSON([]byte(encoded))
	if err != nil {
		t.Error(err)
	}

	if len(testBytes) != p.Size {
		t.Error("incorrect size after unmarshalling")
	}

	if bytes.Compare(testBytes, p.Bytes) != 0 {
		t.Error("bytes before and after marshalling not matching")
	}
}

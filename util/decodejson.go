// Copyright (c) 2018 Kevin Kragenbrink, II
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
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package util

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
)

// Decodejson reads json from an io.Reader (like the body of an http.Request) and converts it to a raw JSON message.
func Decodejson(r io.Reader) (json.RawMessage, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read request body")
	}
	var decoded json.RawMessage
	err = json.Unmarshal(raw, &decoded)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode JSON")
	}
	return decoded, nil
}

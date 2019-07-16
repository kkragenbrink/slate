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

package roll

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoller(t *testing.T) {
	roller, err := NewRoller("cofd", nil)
	roller.SetRand(testRoller)
	assert.Nil(t, err)
	assert.IsType(t, &CofDRollSystem{}, roller)
}

func TestNewRoller_NoRoller(t *testing.T) {
	_, err := NewRoller("test", nil)
	assert.Error(t, err)
}

func TestNewRoller_WithBody(t *testing.T) {
	raw := `{"dice": 5, "again": 10}`
	body := (json.RawMessage)([]byte(raw))
	roller, err := NewRoller("cofd", body)
	assert.Nil(t, err)
	assert.IsType(t, &CofDRollSystem{}, roller)
	assert.Equal(t, int64(10), roller.(*CofDRollSystem).Again)
}

func testRoller(times int, min int64, max int64) []int64 {
	return []int64{}
}

// Copyright (c) 2019 Kevin Kragenbrink, II
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	var got int
	got = Min(2, 0)
	assert.Equal(t, 0, got)
	got = Min(0, 2)
	assert.Equal(t, 0, got)
	got = Min(-1, 2)
	assert.Equal(t, -1, got)
	got = Min(100, 5)
	assert.Equal(t, 5, got)
}

func TestMax(t *testing.T) {
	var got int
	got = Max(2, 0)
	assert.Equal(t, 2, got)
	got = Max(0, 2)
	assert.Equal(t, 2, got)
	got = Max(-1, 2)
	assert.Equal(t, 2, got)
	got = Max(100, 5)
	assert.Equal(t, 100, got)
}

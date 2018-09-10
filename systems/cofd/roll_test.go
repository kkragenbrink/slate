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

package cofd

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func mockRoller(r []int, rr []int) roller {
	rolls := append(r, rr...)

	return func(dice, sides int) []int {
		re := rolls[0:dice]
		rolls = rolls[dice:]
		return re
	}
}

func TestRoll(t *testing.T) {
	o := New()
	es := 1
	er := []int{6, 6, 6, 6, 8}
	err := []int(nil)
	o.SetRoller(mockRoller(er, err))

	s, r, rr := o.roll(5, 10, false, false)
	assert.Equal(t, es, s)
	assert.Equal(t, er, r)
	assert.Equal(t, err, rr)
}

func TestRollAgain(t *testing.T) {
	o := New()
	es := 2
	er := []int{6, 6, 6, 6, 10}
	err := []int{9}
	o.SetRoller(mockRoller(er, err))

	s, r, rr := o.roll(5, 10, false, false)
	assert.Equal(t, es, s)
	assert.Equal(t, er, r)
	assert.Equal(t, err, rr)
}

func TestRollRote(t *testing.T) {
	o := New()
	es := 4
	er := []int{6, 6, 6, 8, 6}
	err := []int{8, 6, 8, 8}
	o.SetRoller(mockRoller(er, err))

	s, r, rr := o.roll(5, 10, true, false)
	assert.Equal(t, es, s)
	assert.Equal(t, er, r)
	assert.Equal(t, err, rr)
}

func TestRollWeak(t *testing.T) {
	o := New()
	es := 3
	er := []int{9, 9, 9, 9, 1}
	err := []int(nil)
	o.SetRoller(mockRoller(er, err))

	s, r, rr := o.roll(5, 10, false, true)
	assert.Equal(t, es, s)
	assert.Equal(t, er, r)
	assert.Equal(t, err, rr)
}

func TestRollChance(t *testing.T) {
	o := New()
	es := 1
	er := []int{10}
	err := []int(nil)
	o.SetRoller(mockRoller(er, err))

	s, r, rr := o.roll(0, 10, false, false)
	assert.Equal(t, es, s)
	assert.Equal(t, er, r)
	assert.Equal(t, err, rr)
}

func TestParseArgs(t *testing.T) {
	argstr := "1 + 2 + 3 + 4 - 10"
	args := strings.Split(argstr, " ")
	expected := 0
	o, e := parseArgs(args)
	assert.Nil(t, e)
	assert.Equal(t, expected, o)
}

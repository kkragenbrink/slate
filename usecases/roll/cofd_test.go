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
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoll(t *testing.T) {
	exs := 1
	exr := []int{6, 6, 6, 6, 8}
	exrr := []int(nil)
	o := genMockCofDRollSystem(mockRoller(exr, exrr), 10, 5, false, false)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(t, err)
	assert.Equal(t, exs, o.Results.Successes)
	assert.Equal(t, exr, o.Results.Rolls)
	assert.Equal(t, exrr, o.Results.Rerolls)
}

func TestRollAgain(t *testing.T) {
	exs := 2
	exr := []int{6, 6, 6, 6, 10}
	exrr := []int{9}
	o := genMockCofDRollSystem(mockRoller(exr, exrr), 10, 5, false, false)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(t, err)
	assert.Equal(t, exs, o.Results.Successes)
	assert.Equal(t, exr, o.Results.Rolls)
	assert.Equal(t, exrr, o.Results.Rerolls)
}

func TestRollRote(t *testing.T) {
	exs := 4
	exr := []int{6, 6, 6, 8, 6}
	exrr := []int{8, 6, 8, 8}
	o := genMockCofDRollSystem(mockRoller(exr, exrr), 10, 5, true, false)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(t, err)
	assert.Equal(t, exs, o.Results.Successes)
	assert.Equal(t, exr, o.Results.Rolls)
	assert.Equal(t, exrr, o.Results.Rerolls)
}

func TestRollWeak(t *testing.T) {
	exs := 3
	exr := []int{9, 9, 9, 9, 1}
	exrr := []int(nil)
	o := genMockCofDRollSystem(mockRoller(exr, exrr), 10, 5, false, true)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(t, err)
	assert.Equal(t, exs, o.Results.Successes)
	assert.Equal(t, exr, o.Results.Rolls)
	assert.Equal(t, exrr, o.Results.Rerolls)
}

func TestRollChance(t *testing.T) {
	exs := 1
	exr := []int{10}
	exrr := []int(nil)
	o := genMockCofDRollSystem(mockRoller(exr, exrr), 10, 5, false, false)
	err := o.Roll(context.Background(), []string{})
	assert.Nil(t, err)
	assert.Equal(t, exs, o.Results.Successes)
	assert.Equal(t, exr, o.Results.Rolls)
	assert.Equal(t, exrr, o.Results.Rerolls)
}

func TestToString(t *testing.T) {
	o := genMockCofDRollSystem(mockRoller([]int{}, []int{}), 9, 4, true, true)
	o.Dice = 4
	o.Verbose = true
	o.Results.Successes = 4
	o.Results.Rolls = []int{8, 8, 9, 7}
	o.Results.Rerolls = []int{8, 4}
	str := o.ToString()
	ex := "rolled 4 CofD Dice (with 9-Again, Rote, Weakness) for 4 Successes. Exceptional success! Rolls: [8 8 9 7] Rerolls: [8 4]"
	assert.Equal(t, ex, str)
}

func mockRoller(r []int, rr []int) roller {
	rolls := append(r, rr...)

	return func(times, min, max int) []int {
		re := rolls[0:times]
		rolls = rolls[times:]
		return re
	}
}

func genMockCofDRollSystem(r roller, again, exceptional int, rote, weak bool) *CofDRollSystem {
	return &CofDRollSystem{
		Again:       again,
		rand:        r,
		Rote:        rote,
		Weakness:    weak,
		Exceptional: exceptional,
	}
}

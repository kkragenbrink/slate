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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CofDTestSuite struct {
	suite.Suite
}

func TestCofD(t *testing.T) {
	suite.Run(t, new(CofDTestSuite))
}

func (suite *CofDTestSuite) TestRoll() {
	exs := 1
	exr := []int64{6, 6, 6, 6, 8}
	exrr := []int64(nil)
	o := genMockCofDRollSystem(cofdMockRoller(exr, exrr), 10, 5, false, false)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), exs, o.Results.Successes)
	assert.Equal(suite.T(), exr, o.Results.Rolls)
	assert.Equal(suite.T(), exrr, o.Results.Rerolls)
}

func (suite *CofDTestSuite) TestRollAgain() {
	exs := 2
	exr := []int64{6, 6, 6, 6, 10}
	exrr := []int64{9}
	o := genMockCofDRollSystem(cofdMockRoller(exr, exrr), 10, 5, false, false)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), exs, o.Results.Successes)
	assert.Equal(suite.T(), exr, o.Results.Rolls)
	assert.Equal(suite.T(), exrr, o.Results.Rerolls)
}

func (suite *CofDTestSuite) TestRollRote() {
	exs := 4
	exr := []int64{6, 6, 6, 8, 6}
	exrr := []int64{8, 6, 8, 8}
	o := genMockCofDRollSystem(cofdMockRoller(exr, exrr), 10, 5, true, false)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), exs, o.Results.Successes)
	assert.Equal(suite.T(), exr, o.Results.Rolls)
	assert.Equal(suite.T(), exrr, o.Results.Rerolls)
}

func (suite *CofDTestSuite) TestRollWeak() {
	exs := 3
	exr := []int64{9, 9, 9, 9, 1}
	exrr := []int64(nil)
	o := genMockCofDRollSystem(cofdMockRoller(exr, exrr), 10, 5, false, true)
	err := o.Roll(context.Background(), []string{"5"})
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), exs, o.Results.Successes)
	assert.Equal(suite.T(), exr, o.Results.Rolls)
	assert.Equal(suite.T(), exrr, o.Results.Rerolls)
}

func (suite *CofDTestSuite) TestRollChance() {
	exs := 1
	exr := []int64{10}
	exrr := []int64(nil)
	o := genMockCofDRollSystem(cofdMockRoller(exr, exrr), 10, 5, false, false)
	var d []string
	err := o.Roll(context.Background(), d)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), exs, o.Results.Successes)
	assert.Equal(suite.T(), exr, o.Results.Rolls)
	assert.Equal(suite.T(), exrr, o.Results.Rerolls)
}

func (suite *CofDTestSuite) TestParseArgs() {
	str := []string{"1+1", "-1", "3-3"}
	exp := 1
	rs := genMockCofDRollSystem(cofdMockRoller([]int64{}, []int64{}), 10, 5, false, false)
	o, err := rs.parseArgs(str)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), exp, o)
}

func (suite *CofDTestSuite) TestToString() {
	o := genMockCofDRollSystem(cofdMockRoller([]int64{}, []int64{}), 9, 4, true, true)
	o.Dice = 4
	o.Verbose = true
	o.Results.Successes = 4
	o.Results.Rolls = []int64{8, 8, 9, 7}
	o.Results.Rerolls = []int64{8, 4}
	str := o.ToString()
	ex := "rolled 4 CofD Dice (with 9-Again, Rote, Weakness) for 4 Successes. Exceptional success! Rolls: [8 8 9 7] Rerolls: [8 4]"
	assert.Equal(suite.T(), ex, str)
}

func cofdMockRoller(r []int64, rr []int64) roller {
	rolls := append(r, rr...)

	return func(times int, min, max int64) ([]int64, error) {
		re := rolls[0:times]
		rolls = rolls[times:]
		return re, nil
	}
}

func genMockCofDRollSystem(r roller, again int64, exceptional int, rote, weak bool) *CofDRollSystem {
	return &CofDRollSystem{
		Again:       again,
		rand:        r,
		Rote:        rote,
		Weakness:    weak,
		Exceptional: exceptional,
	}
}

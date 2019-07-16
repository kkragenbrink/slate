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

package roll

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type D20TestSuite struct {
	suite.Suite
}

func TestD20(t *testing.T) {
	suite.Run(t, new(D20TestSuite))
}

func (suite *D20TestSuite) TestRoll() {
	rolls := []int64{1, 20, 3, 5}
	o := genMockD20RollSystem(d20MockRoller(rolls))
	o.Verbose = true
	err := o.Roll(context.Background(), []string{"4d20"})
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), rolls, o.Expression[0].Rolls)
	assert.Equal(suite.T(), int64(29), o.Expression[0].Value)

	str := o.ToString()
	assert.Equal(suite.T(), "rolled 4d20: 29 (1,20,3,5)", str)
}

func (suite *D20TestSuite) TestRollKeepHighest() {
	rolls := []int64{3, 5, 4, 6}
	o := genMockD20RollSystem(d20MockRoller(rolls))
	err := o.Roll(context.Background(), []string{"4d6kh3"})
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []int64{6, 5, 4, 3}, o.Expression[0].Rolls)
	assert.Equal(suite.T(), int64(15), o.Expression[0].Value)
}

func (suite *D20TestSuite) TestRollKeepLowest() {
	rolls := []int64{3, 5, 4, 6}
	o := genMockD20RollSystem(d20MockRoller(rolls))
	err := o.Roll(context.Background(), []string{"4d6kl3"})
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []int64{3, 4, 5, 6}, o.Expression[0].Rolls)
	assert.Equal(suite.T(), int64(12), o.Expression[0].Value)
}

func d20MockRoller(rolls []int64) roller {
	return func(times int, min, max int64) ([]int64, error) {
		return rolls, nil
	}
}

func genMockD20RollSystem(r roller) *D20RollSystem {
	s := NewD20RollSystem()
	s.SetRand(r)
	return s
}

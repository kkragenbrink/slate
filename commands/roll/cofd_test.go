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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// FormatCofDResults
type FormatCofDResultsSuite struct {
	suite.Suite
	cmd *Command
}

func (s *FormatCofDResultsSuite) SetupTest() {
	cmd := new(Command)
	cmd.flags.system = "cofd"
	cmd.flags.verbose = false
	cmd.flags.cofd.exceptional = 5
	cmd.flags.cofd.again = 10
	cmd.flags.cofd.rote = false
	cmd.flags.cofd.weakness = false
	s.cmd = cmd
}

func (s *FormatCofDResultsSuite) TestNormal() {
	dice := 2
	rolls := []int{10, 7}
	rerolls := []int{8}
	successes := 2

	assert.Equal(s.T(), "rolled 2 CofD dice for 2 successes.", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func (s *FormatCofDResultsSuite) TestVerbose() {
	dice := 2
	rolls := []int{10, 7}
	rerolls := []int{8}
	successes := 2

	s.cmd.flags.verbose = true
	assert.Equal(s.T(), "rolled 2 CofD dice for 2 successes. rolls: [10 7] rerolls: [8]", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func (s *FormatCofDResultsSuite) TestExceptional() {
	dice := 2
	rolls := []int{10, 7}
	rerolls := []int{8}
	successes := 2

	s.cmd.flags.cofd.exceptional = 2
	assert.Equal(s.T(), "rolled 2 CofD dice for 2 successes. Exceptional success!", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func (s *FormatCofDResultsSuite) TestAgain() {
	dice := 2
	rolls := []int{10, 7}
	rerolls := []int{8}
	successes := 2

	s.cmd.flags.cofd.again = 9
	assert.Equal(s.T(), "rolled 2 CofD dice (with 9-again) for 2 successes.", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func (s *FormatCofDResultsSuite) TestRote() {
	dice := 2
	rolls := []int{10, 7}
	rerolls := []int{8}
	successes := 2

	s.cmd.flags.cofd.rote = true
	assert.Equal(s.T(), "rolled 2 CofD dice (with rote) for 2 successes.", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func (s *FormatCofDResultsSuite) TestWeak() {
	dice := 2
	rolls := []int{10, 7}
	rerolls := []int{8}
	successes := 2

	s.cmd.flags.cofd.weakness = true
	assert.Equal(s.T(), "rolled 2 CofD dice (with weakness) for 2 successes.", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func (s *FormatCofDResultsSuite) TestChance() {
	dice := 0
	rolls := []int{5}
	rerolls := []int(nil)
	successes := 0

	assert.Equal(s.T(), "rolled 0 CofD dice (Chance Die) for 0 successes.", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func (s *FormatCofDResultsSuite) TestCriticalFailure() {
	dice := 0
	rolls := []int{1}
	rerolls := []int(nil)
	successes := 0

	assert.Equal(s.T(), "rolled 0 CofD dice (Chance Die) for 0 successes. Critical failure!", s.cmd.formatCofDResults(dice, rolls, rerolls, successes))
}

func TestFormatCofDResult(t *testing.T) {
	suite.Run(t, new(FormatCofDResultsSuite))
}

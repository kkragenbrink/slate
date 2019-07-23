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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize_Basics(t *testing.T) {
	input := "1d20[Basic] + 3[Test] - 5"
	expected := []Token{
		&DiceToken{
			Dice:  1,
			Sides: 20,
			Label: "Basic",
		},
		&ArithmeticToken{
			Multiplier: 1,
		},
		&NumberToken{
			Value: 3,
			Label: "Test",
		},
		&ArithmeticToken{
			Multiplier: -1,
		},
		&NumberToken{
			Value: 5,
		},
	}
	output, err := Tokenize(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, output)
}

func TestTokenize_Fate(t *testing.T) {
	input := "4dF + 3[Lore]"
	expected := []Token{
		&DiceToken{
			Dice:  4,
			Sides: 6,
			Fudge: true,
		},
		&ArithmeticToken{
			Multiplier: 1,
		},
		&NumberToken{
			Value: 3,
			Label: "Lore",
		},
	}
	output, err := Tokenize(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, output)
}

func TestTokenize_Keeps(t *testing.T) {
	input := "4d6kh3 4d6kl3"
	expected := []Token{
		&DiceToken{
			Dice:     4,
			Sides:    6,
			KeepType: KeepTypeHighest,
			Keep:     3,
		},
		&DiceToken{
			Dice:     4,
			Sides:    6,
			KeepType: KeepTypeLowest,
			Keep:     3,
		},
	}
	output, err := Tokenize(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, output)
}

func TestTokenize_Target(t *testing.T) {
	input := "3d10>8 3d10<8 3d10=8 3d10>=8 3d10<=8"
	expected := []Token{
		&DiceToken{
			Dice:       3,
			Sides:      10,
			Comparison: ComparisonGT,
			Target:     8,
		},
		&DiceToken{
			Dice:       3,
			Sides:      10,
			Comparison: ComparisonLT,
			Target:     8,
		},
		&DiceToken{
			Dice:       3,
			Sides:      10,
			Comparison: ComparisonEQ,
			Target:     8,
		},
		&DiceToken{
			Dice:       3,
			Sides:      10,
			Comparison: ComparisonGTE,
			Target:     8,
		},
		&DiceToken{
			Dice:       3,
			Sides:      10,
			Comparison: ComparisonLTE,
			Target:     8,
		},
	}
	output, err := Tokenize(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, output)
}

func TestTokenize_Explode(t *testing.T) {
	input := "3d10! 3d10!>9 3d10>8!>9"
	expected := []Token{
		&DiceToken{
			Dice:    3,
			Sides:   10,
			Explode: 10,
		},
		&DiceToken{
			Dice:    3,
			Sides:   10,
			Explode: 9,
		},
		&DiceToken{
			Dice:       3,
			Sides:      10,
			Comparison: ComparisonGT,
			Target:     8,
			Explode:    9,
		},
	}
	output, err := Tokenize(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, output)
}

func TestTokenize_Failure(t *testing.T) {
	input := "3d10f1 3d10!f1 3d6<4f>5"
	expected := []Token{
		&DiceToken{
			Dice:              3,
			Sides:             10,
			Failure:           1,
			FailureComparison: ComparisonLTE,
		},
		&DiceToken{
			Dice:              3,
			Sides:             10,
			Explode:           10,
			Failure:           1,
			FailureComparison: ComparisonLTE,
		},
		&DiceToken{
			Dice:              3,
			Sides:             6,
			Comparison:        ComparisonLT,
			Target:            4,
			Failure:           5,
			FailureComparison: ComparisonGTE,
		},
	}
	output, err := Tokenize(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, output)
}

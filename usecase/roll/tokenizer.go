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
	"regexp"
	"strconv"
)

// TokenType represents the type of token
type TokenType int

// ComparisonType represents the direction of a comparison
type ComparisonType string

// Constants representing various enums
const (
	KeepTypeHighest int = iota
	KeepTypeLowest

	ComparisonLT  ComparisonType = "<"
	ComparisonGT  ComparisonType = ">"
	ComparisonEQ  ComparisonType = "="
	ComparisonLTE ComparisonType = "<="
	ComparisonGTE ComparisonType = ">="

	diceTokenType TokenType = iota
	numberTokenType
	arithmeticTokenType
)

// The Token interface defines the structure of a roll system token
type Token interface {
	Type() TokenType
}

// A DiceToken represents a random roll
type DiceToken struct {
	Comparison        ComparisonType
	Dice              int
	Explode           int
	Failure           int
	FailureComparison ComparisonType
	Fudge             bool
	Keep              int
	KeepType          int
	Label             string
	Sides             int
	Target            int

	// todo: add support for compounding (!!\d+) and penetrating (!p>\d+) explosions
	// todo: add support for reroll (r\d+)+, (r<\d+)+, (r>\d+)+, and reroll once (ro\d+)+
}

var reDiceToken = regexp.MustCompile(`^\s*(\d+)d(\d+|F)(?:(>|<|=|>=|<=)(\d+))?(!>(\d+)|!)?(?:(f|f>)(\d+))?(?:(kl|kh)(\d+))?(?:\[(.*?)\])?`)

// Type returns the token type for a dice token
func (DiceToken) Type() TokenType { return diceTokenType }

func isDiceToken(b string) bool {
	return reDiceToken.MatchString(b)
}

func parseDiceToken(b string) (*DiceToken, string, error) {
	matches := reDiceToken.FindStringSubmatch(b)
	tok := &DiceToken{}

	// rest
	rest := b[len(matches[0]):]

	// dice
	dice, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, rest, err
	}
	tok.Dice = dice

	// sides
	if matches[2] == "F" {
		tok.Fudge = true
		tok.Sides = 6
	} else {
		sides, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, rest, err
		}
		tok.Sides = sides
	}

	// target
	if matches[3] != "" {
		tok.Comparison = ComparisonType(matches[3])
		target, err := strconv.Atoi(matches[4])
		if err != nil {
			return nil, rest, err
		}
		tok.Target = target
	}

	// explode
	if matches[5] != "" {
		if matches[6] != "" {
			explode, err := strconv.Atoi(matches[6])
			if err != nil {
				return nil, rest, err
			}
			tok.Explode = explode
		} else {
			tok.Explode = tok.Sides
		}
	}

	// weakness/failure
	if matches[8] != "" {
		failure, err := strconv.Atoi(matches[8])
		if err != nil {
			return nil, rest, err
		}
		tok.Failure = failure

		tok.FailureComparison = ComparisonLTE
		if matches[7] == "f>" {
			tok.FailureComparison = ComparisonGTE
		}
	}

	// keep
	if matches[9] != "" {
		if matches[9] == "kh" {
			tok.KeepType = KeepTypeHighest
		}
		if matches[9] == "kl" {
			tok.KeepType = KeepTypeLowest
		}

		keep, err := strconv.Atoi(matches[10])
		if err != nil {
			return nil, rest, err
		}
		tok.Keep = keep
	}

	// label
	if matches[11] != "" {
		tok.Label = matches[11]
	}

	return tok, rest, nil
}

// A NumberToken represents a Numerical value
type NumberToken struct {
	Label string
	Value int
}

var reNumberToken = regexp.MustCompile(`^\s*(\d+)(?:\[(.*?)\])?`)

// Type returns the token type for a number token
func (NumberToken) Type() TokenType { return numberTokenType }

func isNumberToken(b string) bool {
	return reNumberToken.MatchString(b)
}

func parseNumberToken(b string) (*NumberToken, string, error) {
	matches := reNumberToken.FindStringSubmatch(b)
	tok := &NumberToken{}

	// rest
	rest := b[len(matches[0]):]

	// value
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, rest, err
	}
	tok.Value = value

	// label
	if matches[2] != "" {
		tok.Label = matches[2]
	}

	return tok, rest, nil
}

// An ArithmeticToken represents addition or subtraction
type ArithmeticToken struct {
	Multiplier int
}

var reArithmeticToken = regexp.MustCompile(`^\s*([+|-])`)

// Type returns the token type for an arithmetic token
func (ArithmeticToken) Type() TokenType { return arithmeticTokenType }

func isArithmeticToken(b string) bool {
	return reArithmeticToken.MatchString(b)
}

func parseArithmeticToken(b string) (*ArithmeticToken, string, error) {
	matches := reArithmeticToken.FindStringSubmatch(b)
	tok := &ArithmeticToken{}

	// rest
	rest := b[len(matches[0]):]

	// multiplier
	switch matches[1] {
	case "+":
		tok.Multiplier = 1
	case "-":
		tok.Multiplier = -1
	}

	return tok, rest, nil
}

// Tokenize parses an input string into a slice of tokens
func Tokenize(input string) ([]Token, error) {
	tokens := make([]Token, 0)
	for {
		var tok Token
		var err error

		switch {
		case isDiceToken(input):
			tok, input, err = parseDiceToken(input)
		case isNumberToken(input):
			tok, input, err = parseNumberToken(input)
		case isArithmeticToken(input):
			tok, input, err = parseArithmeticToken(input)
		}

		if tok == nil {
			break
		}

		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
	}

	return tokens, nil
}

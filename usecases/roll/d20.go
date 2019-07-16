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
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/kkragenbrink/slate/util"
	"github.com/pkg/errors"
)

const (
	d20Regexp = "([0-9]+)d([0-9]+)(?:(kh|kl)([0-9]+))?"
)

// A D20Token represents a tokenized roll expression
type D20Token struct {
	Dice        int     `json:"dice"`
	Sides       int64   `json:"sides"`
	Value       int64   `json:"value"`
	Rolls       []int64 `json:"rolls"`
	KeepHighest int     `json:"keepHighest"`
	KeepLowest  int     `json:"keepLowest"`
	Negative    bool    `json:"negative"`
}

// The D20RollSystem is the d20 system used in Dungeons & Dragons and similar games.
type D20RollSystem struct {
	rand               roller
	Verbose            bool `json:"verbose"`
	OriginalExpression string
	Expression         []*D20Token `json:"expression"`
	reg                *regexp.Regexp
}

// NewD20RollSystem creates a new instance of the d20 Roll system
func NewD20RollSystem() *D20RollSystem {
	rs := new(D20RollSystem)
	rs.reg, _ = regexp.Compile(d20Regexp)
	return rs
}

// Flags sets up the flag rules for the system
func (rs *D20RollSystem) Flags(fs *flag.FlagSet) {
	fs.BoolVar(&rs.Verbose, "verbose", false, "Whether to use a Verbose output.")

	var system string
	fs.StringVar(&system, "system", "d20", "-- ignored --")
}

// SetRand assigns a random number generator to the system
func (rs *D20RollSystem) SetRand(rand roller) {
	rs.rand = rand
}

// Roll runs the rollsystem for a given set of []tokens.
// This function should only be run once per object.
func (rs *D20RollSystem) Roll(ctx context.Context, tokens []string) error {
	rs.OriginalExpression = strings.Join(tokens, " ")

	if tokens != nil {
		var err error

		rs.Expression, err = rs.parseTokens(tokens)
		if err != nil {
			// todo: wrap probably
			return err
		}
	}

	for _, token := range rs.Expression {
		// If the token has Dice and Sides, we need to roll
		if token.Dice != 0 && token.Sides != 0 {
			rolls, err := rs.rand(token.Dice, 1, token.Sides)
			if err != nil {
				// todo: maybe wrap?
				return err
			}
			token.Rolls = rolls
			if token.KeepHighest != 0 {
				sort.SliceStable(token.Rolls, func(i, j int) bool {
					return token.Rolls[i] > token.Rolls[j]
				})
				for i := 0; i < token.KeepHighest; i++ {
					token.Value += token.Rolls[i]
				}
			}

			if token.KeepLowest != 0 {
				sort.SliceStable(token.Rolls, func(i, j int) bool {
					return token.Rolls[i] < token.Rolls[j]
				})
				for i := 0; i < token.KeepLowest; i++ {
					token.Value += token.Rolls[i]
				}
			}

			if token.KeepLowest == 0 && token.KeepHighest == 0 {
				for i := 0; i < len(token.Rolls); i++ {
					token.Value += token.Rolls[i]
				}
			}
		}

		// If the token is negative, invert the value
		if token.Negative {
			token.Value *= -1
		}
	}

	return nil
}

// ToString converts the Results to a string.
func (rs *D20RollSystem) ToString() string {
	verbose := make([]string, 0)
	var total int64
	for _, token := range rs.Expression {
		total += token.Value

		verb := make([]string, 0)
		for i, roll := range token.Rolls {
			if token.KeepHighest != 0 {
				if i >= token.KeepHighest {
					verb = append(verb, fmt.Sprintf("~~%d~~", roll))
				} else {
					verb = append(verb, fmt.Sprintf("%d", roll))
				}
			}

			if token.KeepLowest != 0 {
				if i >= token.KeepLowest {
					verb = append(verb, fmt.Sprintf("~~%d~~", roll))
				} else {
					verb = append(verb, fmt.Sprintf("%d", roll))
				}
			}

			if token.KeepHighest == 0 && token.KeepLowest == 0 {
				verb = append(verb, fmt.Sprintf("%d", roll))
			}
		}

		verbose = append(verbose, strings.Join(verb, ","))
	}

	if rs.Verbose {
		return fmt.Sprintf("rolled %s: %d (%s)", rs.OriginalExpression, total, strings.Join(verbose, " "))
	}

	return fmt.Sprintf("rolled %s: %d", rs.OriginalExpression, total)
}

func (rs *D20RollSystem) parseTokens(args []string) ([]*D20Token, error) {
	tokens := make([]*D20Token, 0)

	// rejoin all the args so that we can split properly
	argsa := strings.Join(args, "+")

	// ensure negativfes are their own tokens
	argsa = strings.Replace(argsa, "-", "+-+", -1)

	// split on the addition
	args = strings.Split(argsa, "+")

	// reverse the token order so that we can handle subtraction as addition
	gra := util.ReverseStringSlice(args)

	for _, arg := range gra {
		// ignore empty strings
		if arg == "" {
			continue
		}

		// roll?
		roll := rs.reg.FindStringSubmatch(arg)
		if len(roll) != 0 {
			// must be a valid roll!
			dice, _ := strconv.Atoi(roll[1])
			sides, _ := strconv.ParseInt(roll[2], 10, 64)
			keep, _ := strconv.Atoi(roll[4])
			tok := &D20Token{
				Dice:  dice,
				Sides: sides,
			}

			switch roll[3] {
			case "kh":
				tok.KeepHighest = keep
			case "kl":
				tok.KeepLowest = keep
			}

			tokens = append(tokens, tok)
			continue
		}

		// just a number?
		num, err := strconv.Atoi(arg)
		if err != nil {
			// this token is a +, so we can ignore it
			if arg == "+" {
				continue
			}

			// this token is a -, so we need to set the previous token negative
			if arg == "-" {
				x, r := tokens[len(tokens)-1], tokens[:len(tokens)-1]
				x.Negative = true
				r = append(r, x)
				continue
			}

			// this is not a recognized token; bail with an error
			return nil, errors.Wrap(ErrInvalidToken, fmt.Sprintf("token: %s", arg))
		}

		// yup, just a number
		tokens = append(tokens, &D20Token{
			Value: int64(num),
		})
	}

	return tokens, nil
}

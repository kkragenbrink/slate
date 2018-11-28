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
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kkragenbrink/slate/util"
	"strconv"
	"strings"
)

// The CofDRollSystem is the d10 system used in Chronicles of Darkness
type CofDRollSystem struct {
	rand        roller
	verbose     bool
	again       int
	exceptional int
	rote        bool
	weakness    bool
	dice        int
	results     struct {
		successes int
		rolls     []int
		rerolls   []int
	}
}

// Flags sets up the flag rules for the system
func (rs *CofDRollSystem) Flags(fs *flag.FlagSet) {
	fs.BoolVar(&rs.verbose, "verbose", false, "Whether to use a verbose output.")
	fs.IntVar(&rs.again, "again", 10, "The n-again of the rs.")
	fs.IntVar(&rs.exceptional, "exceptional", 5, "The number of successes needed for exceptional success.")
	fs.BoolVar(&rs.rote, "rote", false, "Whether the role is a rote action.")
	fs.BoolVar(&rs.weakness, "weakness", false, "Whether the rs is made with weakness.")
}

// SetRand assigns a random number generator to the system
func (rs *CofDRollSystem) SetRand(rand roller) {
	rs.rand = rand
}

// Roll runs the rollsystem for a given set of []tokens.
// This function should only be run once per object.
func (rs *CofDRollSystem) Roll(ctx context.Context, tokens []string) error {
	var err error
	rs.dice, err = parseArgs(tokens)

	if err != nil {
		// todo: wrap probably
		return err
	}

	successes, rolls, rerolls := rs.doRoll(rs.dice, rs.rote, rs.weakness)
	rs.results.successes = successes
	rs.results.rolls = rolls
	rs.results.rerolls = rerolls

	return nil
}

func (rs *CofDRollSystem) doRoll(dice int, rote, weak bool) (successes int, rolls, rerolls []int) {
	rollsToMake := util.Max(dice, 1)
	var rerollsNeeded int

	rolls = rs.rand(rollsToMake, 1, 10)

	// parse the rolls
	for _, roll := range rolls {
		if dice == 0 {
			if roll == 10 {
				successes++
			}
			continue
		}

		if roll < 8 && rote {
			rerollsNeeded++
			continue
		}

		// success?
		if roll >= 8 {
			successes++
		}

		// reroll?
		if roll >= rs.again {
			rerollsNeeded++
		}

		// weakness?
		if weak && roll == 1 {
			successes--
		}
	}

	// Anything over or equal to the "again" threshhold needs rerolling. Easiest
	// way to do that is just to send it back through.
	if rerollsNeeded > 0 {
		rerollSuccesses, rerollRolls, rerollRerolls := rs.doRoll(rerollsNeeded, false, false)
		successes += rerollSuccesses
		rerolls = append(rerolls, rerollRolls...)
		rerolls = append(rerolls, rerollRerolls...)
	}

	return util.Max(successes, 0), rolls, rerolls
}

// ToJSON converts the results to a JSON object.
// todo: implement toJSON and toString
func (rs *CofDRollSystem) ToJSON() json.RawMessage {
	return nil
}

// ToString converts the results to a string.
func (rs *CofDRollSystem) ToString() string {
	var buff bytes.Buffer
	var flags []string

	buff.WriteString(fmt.Sprintf("rolled %d CofD dice", rs.dice))

	// chance die?
	if rs.dice == 0 {
		buff.WriteString(fmt.Sprintf(" (Chance Die)"))
	}

	// add roll modifying flags
	if rs.again != 10 {
		flags = append(flags, fmt.Sprintf("%d-again", rs.again))
	}
	if rs.rote {
		flags = append(flags, "rote")
	}
	if rs.weakness {
		flags = append(flags, "weakness")
	}
	if len(flags) > 0 {
		buff.WriteString(fmt.Sprintf(" (with %s)", strings.Join(flags, ", ")))
	}

	// add successes
	buff.WriteString(fmt.Sprintf(" for %d successes.", rs.results.successes))

	if rs.results.successes >= rs.exceptional {
		buff.WriteString(" Exceptional success!")
	}

	if rs.dice == 0 && rs.results.rolls[0] == 1 {
		buff.WriteString(" Critical failure!")
	}

	// add rolls and rerolls if desired
	if rs.verbose {
		buff.WriteString(fmt.Sprintf(" rolls: %d", rs.results.rolls))

		if len(rs.results.rerolls) > 0 {
			buff.WriteString(fmt.Sprintf(" rerolls: %d", rs.results.rerolls))
		}
	}

	// send it
	return buff.String()
}

func parseArgs(args []string) (int, error) {
	tokens := []int{0}

	// reverse the token order so that we can handle subtraction as addition
	gra := util.ReverseStringSlice(args)

	// range over the tokens
	for _, arg := range gra {
		num, err := strconv.Atoi(arg)

		if err != nil {
			// this token is a +, so we can ignore it
			if arg == "+" {
				continue
			}

			// this token is a -, so we need to multiply the previous token by -1
			if arg == "-" {
				x, r := tokens[len(tokens)-1], tokens[:len(tokens)-1]
				x = x * -1
				r = append(r, x)
				continue
			}

			// this is not a recognized token; bail with an error
			return 0, err
		}

		// this token is a number; add it to the set
		tokens = append(tokens, num)
	}

	// sum the tokens
	results := 0
	for _, num := range tokens {
		results += num
	}

	return util.Max(results, 0), nil
}

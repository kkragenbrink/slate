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
	"bytes"
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/kkragenbrink/slate/util"
)

// The CofDRollSystem is the d10 system used in Chronicles of Darkness
type CofDRollSystem struct {
	rand        roller
	Verbose     bool  `json:"verbose"`
	Again       int64 `json:"again"`
	Exceptional int   `json:"exceptional"`
	Rote        bool  `json:"rote"`
	Weakness    bool  `json:"weakness"`
	Dice        int   `json:"dice"`
	Results     struct {
		Successes int     `json:"Successes"`
		Rolls     []int64 `json:"Rolls"`
		Rerolls   []int64 `json:"Rerolls"`
	} `json:"Results"`
}

// Flags sets up the flag rules for the system
func (rs *CofDRollSystem) Flags(fs *flag.FlagSet) {
	fs.BoolVar(&rs.Verbose, "verbose", false, "Whether to use a Verbose output.")
	fs.Int64Var(&rs.Again, "again", 10, "The n-Again of the rs.")
	fs.IntVar(&rs.Exceptional, "exceptional", 5, "The number of Successes needed for Exceptional success.")
	fs.BoolVar(&rs.Rote, "rote", false, "Whether the role is a Rote action.")
	fs.BoolVar(&rs.Weakness, "weakness", false, "Whether the rs is made with Weakness.")
}

// SetRand assigns a random number generator to the system
func (rs *CofDRollSystem) SetRand(rand roller) {
	rs.rand = rand
}

// Roll runs the rollsystem for a given set of []tokens.
// This function should only be run once per object.
func (rs *CofDRollSystem) Roll(ctx context.Context, tokens []string) error {
	if tokens != nil {
		var err error

		rs.Dice, err = parseArgs(tokens)
		if err != nil {
			// todo: wrap probably
			return err
		}
	}

	successes, rolls, rerolls, err := rs.doRoll(rs.Dice, rs.Rote, rs.Weakness)
	if err != nil {
		return err
	}
	rs.Results.Successes = successes
	rs.Results.Rolls = rolls
	rs.Results.Rerolls = rerolls

	return nil
}

func (rs *CofDRollSystem) doRoll(dice int, rote, weak bool) (successes int, rolls, rerolls []int64, err error) {
	rollsToMake := util.Max(dice, 1)
	var rerollsNeeded int

	rolls, err = rs.rand(rollsToMake, 1, 10)
	if err != nil {
		return 0, nil, nil, err
	}

	// parse the Rolls
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
		if roll >= rs.Again {
			rerollsNeeded++
		}

		// Weakness?
		if weak && roll == 1 {
			successes--
		}
	}

	// Anything over or equal to the "Again" threshhold needs rerolling. Easiest
	// way to do that is just to send it back through.
	if rerollsNeeded > 0 {
		rerollSuccesses, rerollRolls, rerollRerolls, err := rs.doRoll(rerollsNeeded, false, false)
		if err != nil {
			return 0, nil, nil, err
		}
		successes += rerollSuccesses
		rerolls = append(rerolls, rerollRolls...)
		rerolls = append(rerolls, rerollRerolls...)
	}

	return util.Max(int(successes), 0), rolls, rerolls, err
}

// ToString converts the Results to a string.
func (rs *CofDRollSystem) ToString() string {
	var buff bytes.Buffer
	var flags []string

	buff.WriteString(fmt.Sprintf("rolled %d CofD Dice", rs.Dice))

	// chance die?
	if rs.Dice == 0 {
		buff.WriteString(fmt.Sprintf(" (Chance Die)"))
	}

	// add roll modifying flags
	if rs.Again != 10 {
		flags = append(flags, fmt.Sprintf("%d-Again", rs.Again))
	}
	if rs.Rote {
		flags = append(flags, "Rote")
	}
	if rs.Weakness {
		flags = append(flags, "Weakness")
	}
	if len(flags) > 0 {
		buff.WriteString(fmt.Sprintf(" (with %s)", strings.Join(flags, ", ")))
	}

	// add Successes
	buff.WriteString(fmt.Sprintf(" for %d Successes.", rs.Results.Successes))

	if rs.Results.Successes >= rs.Exceptional {
		buff.WriteString(" Exceptional success!")
	}

	if rs.Dice == 0 && rs.Results.Rolls[0] == 1 {
		buff.WriteString(" Critical failure!")
	}

	// add Rolls and Rerolls if desired
	if rs.Verbose {
		buff.WriteString(fmt.Sprintf(" Rolls: %d", rs.Results.Rolls))

		if len(rs.Results.Rerolls) > 0 {
			buff.WriteString(fmt.Sprintf(" Rerolls: %d", rs.Results.Rerolls))
		}
	}

	// send it
	return buff.String()
}

func parseArgs(args []string) (int, error) {
	tokens := []int{0}

	// rejoin all the args so that we can split properly
	argsa := strings.Join(args, "+")

	// ensure negatives are their own tokens
	argsa = strings.Replace(argsa, "-", "+-+", -1)

	// split on the addition
	args = strings.Split(argsa, "+")

	// reverse the token order so that we can handle subtraction as addition
	gra := util.ReverseStringSlice(args)

	// range over the tokens
	for _, arg := range gra {
		// ignore any empty strings
		if arg == "" {
			continue
		}

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

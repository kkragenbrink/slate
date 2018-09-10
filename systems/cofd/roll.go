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
	"github.com/kkragenbrink/slate/util"
	"strconv"
)

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

// Roll accepts a tokenized dice string, along with the necessary modifiers of
// again, rote, and weakness, then calculates the number of dice to roll and
// calls the system randomizer to determine successes, rolls, and rerolls.
func (c *CofD) Roll(tokens []string, again int, rote, weakness bool) (dice, successes int, rolls, rerolls []int, err error) {
	dice, err = parseArgs(tokens)

	if err != nil {
		return // named return values
	}

	successes, rolls, rerolls = c.roll(dice, again, rote, weakness)
	return // named return values
}

func (c *CofD) roll(dice, again int, rote, weak bool) (successes int, rolls, rerolls []int) {
	rollsToMake := util.Max(dice, 1)
	var rerollsNeeded int

	rolls = c.roller(rollsToMake, 10)

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
		if roll >= again {
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
		rerollSuccesses, rerollRolls, rerollRerolls := c.roll(rerollsNeeded, again, false, false)
		successes += rerollSuccesses
		rerolls = append(rerolls, rerollRolls...)
		rerolls = append(rerolls, rerollRerolls...)
	}

	return util.Max(successes, 0), rolls, rerolls
}

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
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func reverseSlice(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseSlice(input[1:]), input[0])
}

func processCofDArgs(args []string) (int, error) {
	r := []int{0}

	for _, arg := range reverseSlice(args) {
		num, err := strconv.Atoi(arg)

		if err != nil {
			if arg == "+" {
				continue
			}

			if arg == "-" {
				x, r := r[len(r)-1], r[:len(r)-1]
				x = x * -1
				r = append(r, x)
				continue
			}

			return 0, err
		}

		r = append(r, num)
	}

	results := 0
	for _, num := range r {
		results += num
	}
	return max(results, 0), nil
}

func (c *Command) cofd(ctx context.Context, args []string) string {
	dice, err := processCofDArgs(args)
	if err != nil {
		return fmt.Sprintf("%s is not a valid dice algorithm", args)
	}

	rolls, rerolls, successes := cofdroll(dice, c.flags.cofd.again, c.flags.cofd.rote, c.flags.cofd.weakness, rand.Intn)
	return c.formatCofDResults(dice, rolls, rerolls, successes)
}

func (c *Command) formatCofDResults(dice int, rolls []int, rerolls []int, successes int) string {
	// parse the results
	var buff bytes.Buffer
	var flags []string

	// setup response
	buff.WriteString(fmt.Sprintf("rolled %d CofD dice", dice))

	// chance die?
	if dice == 0 {
		buff.WriteString(fmt.Sprintf(" (Chance Die)"))
	}

	// add roll modifying flags
	if c.flags.cofd.again != 10 {
		flags = append(flags, fmt.Sprintf("%d-again", c.flags.cofd.again))
	}
	if c.flags.cofd.rote {
		flags = append(flags, "rote")
	}
	if c.flags.cofd.weakness {
		flags = append(flags, "weakness")
	}
	if len(flags) > 0 {
		buff.WriteString(fmt.Sprintf(" (with %s)", strings.Join(flags, ", ")))
	}

	// add successes
	buff.WriteString(fmt.Sprintf(" for %d successes.", successes))

	if successes >= c.flags.cofd.exceptional {
		buff.WriteString(" Exceptional success!")
	}

	if dice == 0 && rolls[0] == 1 {
		buff.WriteString(" Critical failure!")
	}

	// add rolls and rerolls if desired
	if c.flags.verbose {
		buff.WriteString(fmt.Sprintf(" rolls: %d", rolls))

		if len(rerolls) > 0 {
			buff.WriteString(fmt.Sprintf(" rerolls: %d", rerolls))
		}
	}

	// send it
	return buff.String()
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func cofdroll(dice, again int, rote bool, weak bool, roll func(int) int) (rolls []int, rerolls []int, successes int) {
	rrn := 0
	rtm := dice

	if dice == 0 {
		rtm = 1
	}

	// roll some dice
	for i := 0; i < rtm; i++ {
		die := roll(10) + 1
		rolls = append(rolls, die)

		if dice == 0 {
			if die == 10 {
				successes++
			}
			continue
		}

		if die < 8 && rote {
			die = roll(10) + 1
			rerolls = append(rerolls, die)
		}

		if die >= 8 {
			successes++
		}

		if die >= again {
			rrn++
		}

		if weak && die == 1 {
			successes--
		}
	}

	// handle rerolls
	if rrn > 0 {
		r, rr, s := cofdroll(rrn, again, rote, weak, roll)
		rerolls = append(rerolls, r...)
		rerolls = append(rerolls, rr...)
		successes = successes + s
	}

	return rolls, rerolls, max(successes, 0)
}

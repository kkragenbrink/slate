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
	"flag"
	"math/rand"
)

// Command describes the roll command
type Command struct {
	flags struct {
		system  string
		verbose bool
		cofd    struct {
			again       int
			exceptional int
			rote        bool
			weakness    bool
		}
	}
}

// Name is the name of the command
func (*Command) Name() string { return "roll" }

// Synopsis describes the command's purpose
func (*Command) Synopsis() string { return "Slate's built-in dice roller" }

// Usage describes how to use the command
func (*Command) Usage() string {
	return "roll -system=<system> <dice formula>"
}

// SetFlags sets the flags for the command
func (c *Command) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.flags.system, "system", "cofd", "one of: cofd")
	f.BoolVar(&c.flags.verbose, "verbose", false, "show roll details if true")

	// cofd flags
	f.IntVar(&c.flags.cofd.again, "again", 10, "reroll values which equal or exceed this number")
	f.IntVar(&c.flags.cofd.exceptional, "exceptional", 5, "the number required to get an exceptional success")
	f.BoolVar(&c.flags.cofd.rote, "rote", false, "if true, failures will be rerolled once")
	f.BoolVar(&c.flags.cofd.weakness, "weakness", false, "if true, a roll of 1 will reduce the number of successes")
}

// Execute runs the command
func (c *Command) Execute(ctx context.Context, f *flag.FlagSet) string {
	args := f.Args()

	switch c.flags.system {
	case "cofd":
		return c.cofd(ctx, args)
	}

	return "error processing roll"
}

func (c *Command) roll(sides int) int {
	return rand.Intn(sides) + 1
}

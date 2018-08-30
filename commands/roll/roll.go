//+build !test

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

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
	"encoding/json"
	"errors"
	"flag"
	"github.com/kkragenbrink/slate/infrastructures"
	"github.com/kkragenbrink/slate/usecases"
)

var errInvalidRollSystem = errors.New("choose a roll system (one of: cofd, d20, fate)")

const (
	cofd int = iota
	d20
	fate
)

// A System contains the logic needed to perform a roll.
type System interface {
	Flags(*flag.FlagSet)
	Roll(context.Context, []string) error
	SetRand(roller)
	ToString() string
	ToJSON() json.RawMessage
}

type roller func(times, min, max int) []int

// Roll determines which System to use, parses flags for that system, then hands off processing.
func Roll(ctx context.Context, user *usecases.User, channel *usecases.Channel, fields []string) (System, error) {
	fs := &flag.FlagSet{}
	fs.Usage = func() {}
	var system string
	fs.StringVar(&system, "system", "cofd", "the dice system to use")
	fs.Parse(fields)

	cfs := &flag.FlagSet{}
	rs := NewRoller(system)
	if rs == nil {
		return nil, errInvalidRollSystem
	}
	rs.SetRand(infrastructures.MathRand)
	rs.Flags(cfs)
	cfs.Parse(fields)
	rs.Roll(ctx, cfs.Args())

	return rs, nil
}

// NewRoller creates a new System for the appropriate system.
func NewRoller(system string) System {
	switch system {
	case "cofd":
		return &CofDRollSystem{}
	}

	return nil
}

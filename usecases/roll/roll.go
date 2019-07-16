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
	"flag"

	"github.com/pkg/errors"
)

// ErrInvalidRollSystem is thrown when an invalid roll system is selected
var ErrInvalidRollSystem = errors.New("roll system must be one of: cofd, d20, fate")

// ErrInvalidToken is thrown when an invalid token is sent
var ErrInvalidToken = errors.New("You have submitted an invalid token")

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
}

type roller func(times int, min, max int64) ([]int64, error)

// NewRoller creates a new System for the appropriate system.
func NewRoller(system string, body json.RawMessage) (System, error) {
	var sys System
	switch system {
	case "cofd":
		sys = &CofDRollSystem{}
	case "d20":
		sys = NewD20RollSystem()
	default:
		return nil, ErrInvalidRollSystem
	}

	if body != nil {
		err := json.Unmarshal(body, &sys)
		if err != nil {
			return nil, errors.Wrap(err, "could not decode json")
		}
	}
	return sys, nil
}

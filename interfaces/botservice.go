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

package interfaces

import (
	"context"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/usecases/roll"
	"github.com/pkg/errors"
)

// The BotServiceHandler stores information useful to the bot service message handlers
type BotServiceHandler struct {
}

// A Bot is a repository for discord command handlers
type Bot interface {
	AddHandler(string, BotHandler) error
}

// A BotHandler is a message handler for discord messages
type BotHandler func(ctx context.Context, msg *discordgo.MessageCreate, fields []string) (string, error)

// Roll handles incoming roll messages and sends them to the roll usecase.
func (bs *BotServiceHandler) Roll(ctx context.Context, msg *discordgo.MessageCreate, fields []string) (string, error) {
	// determine the system
	fs := &flag.FlagSet{}
	fs.Usage = func() {}
	var system string
	fs.StringVar(&system, "system", "cofd", "the dice system to use")
	fs.Parse(fields)
	// get a roller
	cfs := &flag.FlagSet{}
	rs, err := roll.NewRoller(system, nil)
	if err != nil {
		// todo: log
		return "", errors.Wrap(err, "could not get a roller")
	}
	rs.Flags(cfs)
	cfs.Parse(fields)
	// roll
	err = rs.Roll(ctx, cfs.Args())
	if err != nil {
		// todo: log
		return "", errors.Wrap(err, "roll failed")
	}
	// send the results
	return rs.ToString(), nil
}

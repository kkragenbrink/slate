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
	"fmt"
	"github.com/kkragenbrink/slate/usecases"
	"github.com/kkragenbrink/slate/usecases/roll"
)

// A Bot is a repository for discord command handlers
type Bot interface {
	AddHandler(string, BotHandler) error
}

// A BotHandler is a handler for a specific discord command
type BotHandler func(ctx context.Context, user *usecases.User, channel *usecases.Channel, fields []string) string

// BotRollHandler handles incoming roll messages and sends them to the roll usecase.
func BotRollHandler(ctx context.Context, user *usecases.User, channel *usecases.Channel, fields []string) string {
	fs := &flag.FlagSet{}
	fs.Usage = func() {}
	var system string
	fs.StringVar(&system, "system", "cofd", "the dice system to use")
	fs.Parse(fields)

	cfs := &flag.FlagSet{}
	rs, err := roll.NewRoller(system, nil)
	if err != nil {
		// todo: log
		return fmt.Sprintf("<@%s> %s", user.ID, err.Error())
	}
	rs.Flags(cfs)
	cfs.Parse(fields)
	err = rs.Roll(ctx, cfs.Args())
	if err != nil {
		// todo: log
		return fmt.Sprintf("<@%s> %s", user.ID, err.Error())
	}
	return fmt.Sprintf("<@%s> %s", user.ID, rs.ToString())
}

func initBotServices(bot Bot) {
	bot.AddHandler("roll", BotRollHandler)
}

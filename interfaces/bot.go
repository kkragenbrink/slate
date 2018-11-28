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
	"fmt"
	"github.com/kkragenbrink/slate/infrastructures"
	"github.com/kkragenbrink/slate/usecases"
	"github.com/kkragenbrink/slate/usecases/roll"
)

// RollHandler handles incoming roll messages and sends them to the roll usecase.
func RollHandler(ctx context.Context, user *usecases.User, channel *usecases.Channel, fields []string) (string, error) {
	roller, err := roll.Roll(ctx, user, channel, fields)
	if err != nil {
		// todo: Wrap probably
		return "", err
	}
	return fmt.Sprintf("<@%s> %s", user.ID, roller.ToString()), nil
}

// Init adds all message handlers to the bot.
func Init(bot *infrastructures.Bot) {
	bot.AddHandler("roll", RollHandler)
}

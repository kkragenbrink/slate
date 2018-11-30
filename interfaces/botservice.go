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
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/domains"
	"github.com/kkragenbrink/slate/interfaces/repositories"
	"github.com/kkragenbrink/slate/usecases/roll"
	"github.com/kkragenbrink/slate/usecases/sheet"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// SiteURL is the URL of the website which serves static content
// todo: this should be a configuration parameter
const SiteURL = "https://slate.sosly.org"

// The BotServiceHandler stores information useful to the bot service message handlers
type BotServiceHandler struct {
	bot Bot
	db  repositories.Database
}

// NewBotServiceHandler creates a new BotServiceHandler instance
func NewBotServiceHandler(bot Bot, db repositories.Database) *BotServiceHandler {
	bsh := new(BotServiceHandler)
	bsh.bot = bot
	bsh.db = db
	return bsh
}

// A Bot is a repository for discord command handlers
type Bot interface {
	AddHandler(string, BotHandler) error
	Channel(string) (*discordgo.Channel, error)
	User(string) (*discordgo.User, error)
}

// A BotHandler is a message handler for discord messages
type BotHandler func(ctx context.Context, msg *discordgo.MessageCreate, fields []string) (string, error)

// Sheet creates a new character sheet
func (bs *BotServiceHandler) Sheet(ctx context.Context, msg *discordgo.MessageCreate, fields []string) (string, error) {
	// determine the sheet system
	fs := &flag.FlagSet{}
	fs.Usage = func() {}
	var system string
	fs.StringVar(&system, "system", "wtf2e", "the character system to use")
	fs.Parse(fields)
	name := strings.Join(fs.Args(), " ")
	player, _ := strconv.ParseInt(msg.Author.ID, 10, 64)
	ch, err := bs.bot.Channel(msg.ChannelID)
	if err != nil {
		// todo
		return "", err
	}
	guild, _ := strconv.ParseInt(ch.GuildID, 10, 64)
	repo := bs.db.Repository("character").(domains.CharacterRepository)
	character, err := sheet.New(ctx, repo, name, system, guild, player)
	if err != nil {
		return "", errors.Wrap(err, "could not create a new sheet")
	}
	return fmt.Sprintf("your new character is at %s/sheets/%s", SiteURL, character.ID), nil
}

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

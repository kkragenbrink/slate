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

package infrastructures

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/settings"
	"github.com/kkragenbrink/slate/usecases"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

var errDuplicateHandler = errors.New("duplicate handler already exists")

// A DiscordSession contains instructions for communicating with discord.
type DiscordSession interface {
	AddHandler(handler interface{}) func()
	Channel(string) (*discordgo.Channel, error)
	ChannelMessageSend(string, string) (*discordgo.Message, error)
	Close() error
	Guild(string) (*discordgo.Guild, error)
	Open() error
}

// The Bot contains the connection to discord as well as the message handlers.
type Bot struct {
	settings *settings.Settings
	handlers []*MessageHandler
	session  DiscordSession
	mutex    *sync.Mutex
}

// A MessageHandler is a command name and a handler function
type MessageHandler struct {
	command string
	handle  MessageHandlerFunc
}

// A MessageHandlerFunc is a context-aware response to a message from Discord.
type MessageHandlerFunc func(ctx context.Context, user *usecases.User, channel *usecases.Channel, fields []string) (string, error)

// NewBot returns a new Discord Bot, which will be used by the application for
// communicating with discord.
func NewBot(set *settings.Settings) (*Bot, error) {
	connstr := fmt.Sprintf("Bot %s", set.DiscordToken)

	session, err := discordgo.New(connstr)
	if err != nil {
		return nil, errors.Wrap(err, "could not establish discord session")
	}

	bot := new(Bot)
	bot.settings = set
	bot.session = session
	bot.mutex = &sync.Mutex{}

	return bot, nil
}

// AddHandler adds a new MessageHandler to the bot.
func (bot *Bot) AddHandler(command string, handle MessageHandlerFunc) error {
	if bot.hasHandler(command) {
		return errDuplicateHandler
	}
	handler := new(MessageHandler)
	handler.command = command
	handler.handle = handle
	bot.handlers = append(bot.handlers, handler)
	return nil
}

func (bot *Bot) hasHandler(command string) bool {
	for _, handler := range bot.handlers {
		if handler.command == command {
			return true
		}
	}
	return false
}

func (bot *Bot) handleMessageCreateInterface(session *discordgo.Session, msg *discordgo.MessageCreate) {
	bot.handleMessageCreate(session, msg)
}

func (bot *Bot) handleMessageCreate(session DiscordSession, msg *discordgo.MessageCreate) {
	if !strings.HasPrefix(msg.Content, bot.settings.CommandPrefix) {
		return // this is not our command to handle
	}

	cmd := strings.Trim(msg.Content[len(bot.settings.CommandPrefix):], " ")
	name := strings.Fields(cmd)[0]
	fields := strings.Fields(cmd)[1:]

	if !bot.hasHandler(name) {
		return // still not our command to handle
	}

	// todo: this should be a setting instead of a magic number
	timeout, err := time.ParseDuration("5s")
	if err != nil {
		bot.session.ChannelMessageSend(msg.ChannelID, "Woops, Something went wrong!")
		return
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	bot.mutex.Lock()
	defer bot.mutex.Unlock()
	for _, handler := range bot.handlers {
		if handler.command == name {
			// get channel details
			dc, err := bot.session.Channel(msg.ChannelID)
			if err != nil {
				// todo: log error
				bot.session.ChannelMessageSend(msg.ChannelID, "Woops, Something went wrong!")
				return
			}
			channel := usecases.NewChannel(dc.ID, dc.GuildID)

			// get user details
			user := usecases.NewUser(msg.Author.ID, msg.Author.Username)

			// handle it
			response, err := handler.handle(ctx, user, channel, fields)
			if err != nil {
				// todo: log error
				bot.session.ChannelMessageSend(msg.ChannelID, err.Error())
				return
			}
			bot.session.ChannelMessageSend(msg.ChannelID, response)

			return
		}
	}
}

// Start establishes the connection to discord and adds the message handler.
func (bot *Bot) Start() error {
	err := bot.session.Open()
	if err != nil {
		return errors.Wrap(err, "could not connect to discord")
	}
	bot.session.AddHandler(bot.handleMessageCreateInterface)

	return nil
}

// Stop closes the connection to discord.
func (bot *Bot) Stop() error {
	err := bot.session.Close()
	if err != nil {
		return errors.Wrap(err, "error disconnecting from discord")
	}
	return nil
}

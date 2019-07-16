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

package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/interfaces"
	"github.com/kkragenbrink/slate/settings"
	"github.com/pkg/errors"
)

var errDuplicateHandler = errors.New("duplicate handler already exists")

// A DiscordSession contains instructions for communicating with discord.
type DiscordSession interface {
	AddHandler(handler interface{}) func()
	Channel(string) (*discordgo.Channel, error)
	ChannelMessageSend(string, string) (*discordgo.Message, error)
	Close() error
	GuildChannels(string) ([]*discordgo.Channel, error)
	Open() error
	User(string) (*discordgo.User, error)
}

// The Bot contains the connection to discord as well as the message handlers.
type Bot struct {
	logger     *SlateLogger
	settings   *settings.Settings
	handlers   []*BotMessageHandler
	svchandler *interfaces.BotServiceHandler
	session    DiscordSession
	mutex      *sync.Mutex
}

// A BotMessageHandler is a command name and a handler function
type BotMessageHandler struct {
	command string
	handle  interfaces.BotHandler
}

// NewBot returns a new Discord bot, which will be used by the application for
// communicating with discord.
func NewBot(set *settings.Settings, db *DatabaseService, rand *RandomService) (*Bot, error) {
	connstr := fmt.Sprintf("Bot %s", set.DiscordToken)

	session, err := discordgo.New(connstr)
	if err != nil {
		return nil, errors.Wrap(err, "could not establish discord session")
	}

	bot := new(Bot)
	bot.logger = NewSlateLogger()
	bot.initServiceHandler(db, rand)
	bot.settings = set
	bot.session = session
	bot.mutex = &sync.Mutex{}

	return bot, nil
}

func (bot *Bot) initServiceHandler(db *DatabaseService, rand *RandomService) {
	bs := interfaces.NewBotServiceHandler(bot, db, rand)
	bot.AddHandler("sheet", bs.Sheet)
	bot.AddHandler("roll", bs.Roll)
	bot.svchandler = bs
}

// AddHandler adds a new BotHandler to the bot.
func (bot *Bot) AddHandler(command string, handle interfaces.BotHandler) error {
	if bot.hasHandler(command) {
		return errDuplicateHandler
	}
	handler := new(BotMessageHandler)
	handler.command = command
	handler.handle = handle
	bot.handlers = append(bot.handlers, handler)
	return nil
}

// Channel gets a channel object by ID
func (bot *Bot) Channel(id string) (*discordgo.Channel, error) {
	ch, err := bot.session.Channel(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not find channel by id")
	}
	return ch, nil
}

// Channels gets the list of channel objects for a guild
func (bot *Bot) Channels(id string) ([]*discordgo.Channel, error) {
	fmt.Println("id =", id)
	channels := make([]*discordgo.Channel, 0)
	ch, err := bot.session.GuildChannels(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not find guild by id")
	}
	for _, channel := range ch {
		if channel.Type == discordgo.ChannelTypeGuildText {
			channels = append(channels, channel)
		}
	}
	return channels, nil
}

// SendMessage sends a message to a specified channel
func (bot *Bot) SendMessage(id string, message string) error {
	_, err := bot.session.ChannelMessageSend(id, message)
	return err
}

// User gets a user object by ID
func (bot *Bot) User(id string) (*discordgo.User, error) {
	u, err := bot.session.User(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not find user by id")
	}
	return u, nil
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
	start := time.Now()
	if !strings.HasPrefix(msg.Content, bot.settings.CommandPrefix) {
		return // this is not our command to handle
	}

	cmd := strings.Trim(msg.Content[len(bot.settings.CommandPrefix):], " ")
	name := strings.Fields(cmd)[0]
	fields := strings.Fields(cmd)[1:]

	if !bot.hasHandler(name) {
		return // still not our command to handle
	}
	log := bot.logger.NewDiscordLogEntry(msg, name, fields)

	// todo: this should be a setting instead of a magic number
	timeout, err := time.ParseDuration("5s")
	if err != nil {
		bot.session.ChannelMessageSend(msg.ChannelID, "Woops, Something went wrong!")
		log.Write(LogError, 0, time.Since(start))
		return
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	bot.mutex.Lock()
	defer bot.mutex.Unlock()
	for _, handler := range bot.handlers {
		if handler.command == name {
			// handle it
			results, err := handler.handle(ctx, msg, fields)
			if err != nil {
				response := fmt.Sprintf("%s %s", msg.Author.Mention(), err.Error())
				bot.session.ChannelMessageSend(msg.ChannelID, response)
				log.Write(LogWarn, len(response), time.Since(start))
			} else {
				response := fmt.Sprintf("%s %s", msg.Author.Mention(), results)
				bot.session.ChannelMessageSend(msg.ChannelID, response)
				log.Write(LogInfo, len(response), time.Since(start))
			}
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

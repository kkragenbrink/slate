// Copyright (c) 2019 Kevin Kragenbrink, II
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
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/domain"
	"github.com/sirupsen/logrus"
)

// A DiscordSession contains instructions for communicating with Discord
type DiscordSession interface {
	AddHandler(handler interface{}) func()
	Channel(id string) (*discordgo.Channel, error)
	ChannelMessageSend(id string, message string) (*discordgo.Message, error)
	Close() error
	User(id string) (*discordgo.User, error)
}

// A BotMessageCreateHandler is a method for handling a message from Discord
type BotMessageCreateHandler func(ctx context.Context, msg *discordgo.MessageCreate, fields []string) (string, error)

// The Bot contains the connection to discord, as well as message handlers
type Bot struct {
	Config   *domain.SlateConfig
	Handlers map[string]BotMessageCreateHandler
	Mutex    *sync.Mutex
	Random   *Random
	Session  *discordgo.Session
}

// NewBot instantiates a new Bot
func NewBot(config *domain.SlateConfig, random *Random) (*Bot, error) {
	connStr := fmt.Sprintf("Bot %s", config.DiscordToken)

	session, err := discordgo.New(connStr)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Config:   config,
		Handlers: make(map[string]BotMessageCreateHandler),
		Mutex:    &sync.Mutex{},
		Random:   random,
		Session:  session,
	}

	bot.RegisterRoutes()

	return bot, nil
}

// AddHandler adds a new handler to the bot
func (bot *Bot) AddHandler(command string, handle BotMessageCreateHandler) {
	bot.Handlers[command] = handle
}

func (bot *Bot) handleMessageCreateInterface(session *discordgo.Session, msg *discordgo.MessageCreate) {
	bot.HandleMessageCreate(msg)
}

// HandleMessageCreate responds to MessageCreate events from Discord and passes them to an appropriate handler
func (bot *Bot) HandleMessageCreate(msg *discordgo.MessageCreate) {
	logrus.Debug(msg.Content)
	if !strings.HasPrefix(msg.Content, bot.Config.BotCommandPrefix) {
		return // not our command
	}

	parsedInput := strings.Trim(msg.Content[len(bot.Config.BotCommandPrefix):], " ")
	command := strings.Fields(parsedInput)[0]
	fields := strings.Fields(parsedInput)[1:]

	var response string
	var err error

	if !bot.HasHandler(command) {
		logrus.WithField("command", command).Warn("no handler found")
		return // we don't have a handler for this command
	}

	start := time.Now()
	defer func() {
		lf := logrus.Fields{
			"command": command,
			"took":    time.Since(start),
		}
		if err != nil {
			lf["err"] = err
		} else {
			lf["size"] = len([]byte(response))
		}

		logrus.WithFields(lf).Info("message sent")
	}()

	ctx, cancelFunc := context.WithTimeout(context.Background(), bot.Config.BotCommandTimeout)
	defer cancelFunc()

	bot.Mutex.Lock()
	defer bot.Mutex.Unlock()
	handle := bot.Handlers[command]
	response, err = handle(ctx, msg, fields)

	if err != nil {
		response = err.Error()
	}
	bot.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s %s", msg.Author.Mention(), response))
}

// HasHandler checks whether the bot is able to handle the given command
func (bot *Bot) HasHandler(command string) bool {
	_, ok := bot.Handlers[command]
	return ok
}

// RegisterRoutes binds the routes for the discord bot to the handler methods
func (bot *Bot) RegisterRoutes() {
	bot.AddHandler("roll", bot.handleRoll)
}

// Start establishes the connection to discord and begins watching for messages
func (bot *Bot) Start() error {
	err := bot.Session.Open()
	if err != nil {
		return err
	}
	bot.Session.AddHandler(bot.handleMessageCreateInterface)
	return nil
}

// Stop closes the connection to discord
func (bot *Bot) Stop() error {
	return bot.Session.Close()
}

func (bot *Bot) handleRoll(ctx context.Context, msg *discordgo.MessageCreate, fields []string) (string, error) {
	
}

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

package bot

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/config"
	"strings"
	"time"
)

// Bot contains the connection information and commands used to communicate with
// Discord.
type Bot struct {
	config   *config.Config
	session  DiscordSession
	commands []SlateCommand
}

// AddCommand allows a new command to be registered with the bot
func (b *Bot) AddCommand(command SlateCommand) {
	b.commands = append(b.commands, command)
}

// Commands returns the list of commands
func (b *Bot) Commands() []SlateCommand {
	return b.commands
}

// Stop cleanly closes the connection to Discord
func (b *Bot) Stop() error {
	err := b.session.Close()
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleMessageCreate(s DiscordSession, msg *discordgo.MessageCreate) {
	if !strings.HasPrefix(msg.Content, b.config.CommandPrefix) {
		return // this is not our command to handle
	}

	timeout, err := time.ParseDuration("5s")
	if err != nil {
		b.session.ChannelMessageSend(msg.ChannelID, "Woops, Something went wrong!")
		return
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	cmd := strings.Trim(msg.Content[len(b.config.CommandPrefix):], " ")
	name := strings.Fields(cmd)[0]
	fields := strings.Fields(cmd)[1:]

	var commandList [][]string
	var longest int

	for _, command := range b.commands {
		if name == command.Name() {
			command.Execute(ctx, fields, b.session, msg)
			return
		}

		length := len(command.Name())
		if length > longest {
			longest = length
		}
		usage := []string{command.Name(), command.Synopsis()}
		commandList = append(commandList, usage)
	}

	// command not found
	lines := []string{fmt.Sprintf("\n```Usage: %s<command> <args>\n\nCommands:", b.config.CommandPrefix)}
	template := fmt.Sprintf("\t%%%ds\t%%s", longest)
	for _, command := range commandList {
		lines = append(lines, fmt.Sprintf(template, command[0], command[1]))
	}
	lines = append(lines, "```")
	b.session.ChannelMessageSend(msg.ChannelID, strings.Join(lines, "\n"))
}

// New returns a new instance of a Bot and establishes the connection to
// discord.
func New(cfg *config.Config, s DiscordSession) (*Bot, error) {
	b := new(Bot)
	b.config = cfg
	b.session = s

	// open the websocket
	err := s.Open()
	if err != nil {
		return nil, err
	}

	// establish the command handler
	s.AddHandler(b.handleMessageCreate)

	return b, nil
}

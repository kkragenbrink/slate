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
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/config"
)

// Bot contains the connection information and commands used to communicate with
// Discord.
type Bot struct {
	config   *config.Config
	session  DiscordSession
	commands []SlateCommand
	mutex    *sync.Mutex
}

// AddCommand allows a new command to be registered with the bot
func (b *Bot) AddCommand(command SlateCommand) {
	b.mutex.Lock()
	b.commands = append(b.commands, command)
	b.mutex.Unlock()
}

// Commands returns the list of commands
func (b *Bot) Commands() []SlateCommand {
	b.mutex.Lock()
	defer b.mutex.Unlock()
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

func (b *Bot) handleMessageCreate(msg *discordgo.MessageCreate) {
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

	b.mutex.Lock()
	defer b.mutex.Unlock()
	for _, command := range b.commands {
		if name == command.Name() {
			fmt.Printf("user: %s#%s, command: %s\n", msg.Author.Username, msg.Author.Discriminator, command.Name())
			fs := &flag.FlagSet{}
			command.SetFlags(fs)
			fs.Parse(fields)

			output := command.Execute(ctx, fs)
			response := fmt.Sprintf("%s %s", msg.Author.Mention(), output)
			b.session.ChannelMessageSend(msg.ChannelID, response)
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
	b := &Bot{
		config:  cfg,
		session: s,
		mutex:   &sync.Mutex{},
	}

	// open the websocket
	err := s.Open()
	if err != nil {
		return nil, err
	}

	// establish the command handler
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		b.handleMessageCreate(m)
	})

	return b, nil
}

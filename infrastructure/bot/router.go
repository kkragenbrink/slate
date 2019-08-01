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

package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// DiscordSession is an interface for mocking a discord session
type DiscordSession interface {
	ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
}

// Handle is a handle for discord Message types
type Handle func(m *discordgo.Message) string

// MessageType is a type of Discord message
type MessageType int

// Constants representing MessageType
const (
	MessageCreate MessageType = iota
	MessageUpdate
)

// Router is a means of routing messages to handlers
type Router struct {
	events map[MessageType]map[string]Handle
	prefix string
}

// NewRouter creates a new Discord router
func NewRouter() *Router {
	events := make(map[MessageType]map[string]Handle)
	return &Router{
		events: events,
	}
}

// AddCommand adds a MessageCreate handler
func (r *Router) AddCommand(command string, handle Handle) {
	if _, ok := r.events[MessageCreate]; !ok {
		r.events[MessageCreate] = make(map[string]Handle)
	}

	r.events[MessageCreate][command] = handle
}

// RouteMessageCreate is assigned to a discordgo session as a message handler for create messages
func (r *Router) RouteMessageCreate(s DiscordSession, m *discordgo.MessageCreate) {
	// incorrect command prefix
	if r.prefix != "" && !strings.HasPrefix(m.Content, r.prefix) {
		logrus.WithField("command", m.Content).Debug("unhandled command")
		return
	}

	if _, ok := r.events[MessageCreate]; !ok {
		logrus.Error("no commands registered for bot")
		s.ChannelMessageSend(m.ChannelID, "no commands registered for bot")
		return
	}

	cmdline := strings.Trim(m.Content[len(r.prefix):], " ")
	command := strings.Fields(cmdline)[0]

	handle, ok := r.events[MessageCreate][command]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("command `%s` not found", command))
		return
	}

	go func(handle Handle) {
		response := handle(m.Message)
		s.ChannelMessageSend(m.ChannelID, response)
	}(handle)
}

// SetPrefix sets the command prefix for messages from discord
func (r *Router) SetPrefix(p string) {
	r.prefix = p
}

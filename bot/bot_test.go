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
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNew tests the New() function without any arguments.  This should return a
// valid Bot{} struct and no errors.
func TestNew(t *testing.T) {
	cfg := new(config.Config)
	cfg.DiscordToken = "test-token"

	md := new(MockDiscordSession)
	md.On("Open").Return(nil)
	md.On("AddHandler").Return(nil)

	_, err := New(cfg, md)
	assert.Nil(t, err)

	md.AssertExpectations(t)
}

// TestStop tests the Stop() function.
func TestStop(t *testing.T) {
	b := new(Bot)

	md := new(MockDiscordSession)
	md.On("Close").Return(nil)
	b.session = md

	err := b.Stop()
	assert.Nil(t, err)

	md.AssertExpectations(t)
}

// TestHandleMessageCreate tests the ability to handle an incoming message
func TestHandleMessageCreate(t *testing.T) {
	b := new(Bot)

	cfg := new(config.Config)
	cfg.CommandPrefix = "$test"
	b.config = cfg

	md := new(MockDiscordSession)
	md.On("ChannelMessageSend").Return(nil)
	b.session = md

	mc := new(MockCommand)
	b.AddCommand(mc)

	msg := &discordgo.Message{
		ChannelID: "test",
		Content:   "$test test",
	}
	msgc := &discordgo.MessageCreate{Message: msg}

	b.handleMessageCreate(md, msgc)

	md.AssertExpectations(t)
	mc.AssertExpectations(t)
}

func Test_Commands(t *testing.T) {
	b := new(Bot)

	cfg := new(config.Config)
	cfg.CommandPrefix = "$test"
	b.config = cfg

	md := new(MockDiscordSession)
	b.session = md

	mc := new(MockCommand)
	mc.On("Execute").Return(nil)

	b.AddCommand(mc)
	msg := &discordgo.Message{
		ChannelID: "test",
		Content:   "$test mock",
	}
	msgc := &discordgo.MessageCreate{Message: msg}

	b.handleMessageCreate(md, msgc)
	md.AssertExpectations(t)
	mc.AssertExpectations(t)
}

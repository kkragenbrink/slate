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
	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/kkragenbrink/slate/infrastructures/mocks"
	"github.com/kkragenbrink/slate/usecases"
	"github.com/pkg/errors"
	"testing"

	"github.com/kkragenbrink/slate/settings"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := new(settings.Settings)
	cfg.DiscordToken = "test-token"
	_, err := NewBot(cfg)
	assert.Nil(t, err)
}

func TestAddHandler(t *testing.T) {
	bot := new(Bot)

	// positive flow
	bot.AddHandler("test", genMockHandler("test"))
	assert.Equal(t, 1, len(bot.handlers))
	assert.True(t, bot.hasHandler("test"))

	// negative flow
	err := bot.AddHandler("test", genMockHandler("test"))
	assert.Error(t, errDuplicateHandler, err)
}

func TestHandleMessageCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bot, _ := NewBot(&settings.Settings{CommandPrefix: "$"})
	session := mocks.NewMockDiscordSession(ctrl)
	author := "a1"
	channel := "c1"
	guild := "g1"
	msg := "$test"
	message := genMockMessage(author, channel, msg)
	handler := genMockHandler("test")
	mockChannel := genMockChannel(channel, guild)
	bot.AddHandler("test", handler)
	session.EXPECT().Channel(gomock.Eq(channel)).Return(mockChannel, nil)
	session.EXPECT().ChannelMessageSend(gomock.Eq(channel), gomock.Eq("test"))
	bot.session = session
	bot.handleMessageCreate(session, message)
}

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Open()
	session.EXPECT().AddHandler(gomock.Any())
	bot.session = session
	bot.Start()
}

func TestStartError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Open().Return(errSampleError)
	bot.session = session
	err := bot.Start()
	assert.Error(t, errSampleError, err)
}

func TestStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Close()

	bot.session = session
	bot.Stop()
}

func TestStopError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Close().Return(errSampleError)
	bot.session = session
	err := bot.Stop()
	assert.Error(t, errSampleError, err)
}

func genMockChannel(channel, guild string) *discordgo.Channel {
	return &discordgo.Channel{ID: channel, GuildID: guild}
}

func genMockHandler(response string) MessageHandlerFunc {
	return func(ctx context.Context, u *usecases.User, c *usecases.Channel, s []string) (string, error) {
		return response, nil
	}
}

func genMockMessage(author, channel, msg string) *discordgo.MessageCreate {
	message := &discordgo.Message{}
	message.ChannelID = channel
	message.Author = &discordgo.User{}
	message.Author.ID = author
	message.Author.Username = author + " name"
	message.Content = msg
	msgc := &discordgo.MessageCreate{Message: message}
	return msgc
}

var errSampleError = errors.New("sample")

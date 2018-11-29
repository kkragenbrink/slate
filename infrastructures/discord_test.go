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
	"github.com/kkragenbrink/slate/interfaces"
	"github.com/kkragenbrink/slate/settings"
	"github.com/kkragenbrink/slate/usecases"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BotSuite struct {
	suite.Suite
}

func TestBotSuite(t *testing.T) {
	suite.Run(t, new(BotSuite))
}

func (suite *BotSuite) TestNewBot() {
	set := new(settings.Settings)
	set.DiscordToken = "test-token"
	_, err := NewBot(set)
	assert.Nil(suite.T(), err)
}

func (suite *BotSuite) TestAddHandler() {
	bot := new(Bot)

	// positive flow
	bot.AddHandler("test", genMockHandler("test"))
	assert.Equal(suite.T(), 1, len(bot.handlers))
	assert.True(suite.T(), bot.hasHandler("test"))

	// negative flow
	err := bot.AddHandler("test", genMockHandler("test"))
	assert.Error(suite.T(), errDuplicateHandler, err)
}

func (suite *BotSuite) TestHandleMessageCreate() {
	ctrl := gomock.NewController(suite.T())
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

func (suite *BotSuite) TestStart() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Open()
	session.EXPECT().AddHandler(gomock.Any())
	bot.session = session
	bot.Start()
}

func (suite *BotSuite) TestStartError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Open().Return(errSampleError)
	bot.session = session
	err := bot.Start()
	assert.Error(suite.T(), errSampleError, err)
}

func (suite *BotSuite) TestStop() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Close()

	bot.session = session
	bot.Stop()
}

func (suite *BotSuite) TestStopError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	bot := new(Bot)
	session := mocks.NewMockDiscordSession(ctrl)
	session.EXPECT().Close().Return(errSampleError)
	bot.session = session
	err := bot.Stop()
	assert.Error(suite.T(), errSampleError, err)
}

func genMockChannel(channel, guild string) *discordgo.Channel {
	return &discordgo.Channel{ID: channel, GuildID: guild}
}

func genMockHandler(response string) interfaces.BotHandler {
	return func(ctx context.Context, u *usecases.User, c *usecases.Channel, s []string) string {
		return response
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

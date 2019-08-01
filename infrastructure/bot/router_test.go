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
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RouterSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	sess   *MockDiscordSession
	router *Router
}

func (s *RouterSuite) BeforeTest(suiteName, testName string) {
	s.ctrl = gomock.NewController(s.T())
	s.sess = NewMockDiscordSession(s.ctrl)
	s.router = NewRouter()
}

func (s *RouterSuite) AfterTest(suiteName, testName string) {
	s.ctrl.Finish()
}

func (s *RouterSuite) TestAddCommand() {
	command := "test"
	expected := "expected"
	handleFunc := func(m *discordgo.Message) string { return expected }

	s.router.AddCommand(command, handleFunc)
	assert.Equal(s.T(), 1, len(s.router.events[MessageCreate]))
	assert.NotNil(s.T(), s.router.events[MessageCreate][command])

	output := s.router.events[MessageCreate][command](&discordgo.Message{})
	assert.Equal(s.T(), "expected", output)
}

func (s *RouterSuite) TestCommand() {
	channelID := "channel-test"
	command := "test"
	input := "test"
	expected := "expected"
	wg := &sync.WaitGroup{}
	handleFunc := func(m *discordgo.Message) string { return expected }
	msg := newMessageCreate(channelID, "user", input)
	s.router.AddCommand(command, handleFunc)
	s.sess.EXPECT().ChannelMessageSend(channelID, expected).Do(func(a, b string) { wg.Done() })
	wg.Add(1)
	s.router.RouteMessageCreate(s.sess, msg)
	wg.Wait()
}

func (s *RouterSuite) TestCommandWithPrefix() {
	channelID := "channel-test"
	command := "test"
	input := "$test"
	input1 := "test"
	expected := "expected"
	wg := &sync.WaitGroup{}
	handleFunc := func(m *discordgo.Message) string { return expected }
	msg := newMessageCreate(channelID, "user", input)
	msg1 := newMessageCreate(channelID, "user", input1)

	s.router.SetPrefix("$")
	s.router.AddCommand(command, handleFunc)

	s.sess.EXPECT().ChannelMessageSend(channelID, expected).Do(func(a, b string) { wg.Done() })

	wg.Add(1)
	s.router.RouteMessageCreate(s.sess, msg)
	s.router.RouteMessageCreate(s.sess, msg1)
	wg.Wait()
}

func (s *RouterSuite) TestCommandNotFound() {
	channelID := "channel-test"
	command := "test"
	input := "testing"
	expected := "command `testing` not found"
	handleFunc := func(m *discordgo.Message) string { return expected }
	msg := newMessageCreate(channelID, "user", input)
	s.router.AddCommand(command, handleFunc)
	s.sess.EXPECT().ChannelMessageSend(channelID, expected)
	s.router.RouteMessageCreate(s.sess, msg)
}

func (s *RouterSuite) TestNoCommands() {
	channelID := "channel-test"
	input := "testing"
	expected := "no commands registered for bot"
	msg := newMessageCreate(channelID, "user", input)
	s.sess.EXPECT().ChannelMessageSend(channelID, expected)
	s.router.RouteMessageCreate(s.sess, msg)
}

func newMessageCreate(channel, user, message string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID:       user,
				Username: user,
			},
			ChannelID: channel,
			Content:   message,
		},
	}
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}

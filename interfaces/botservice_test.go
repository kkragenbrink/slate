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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/kkragenbrink/slate/infrastructure/bot"

	"github.com/kkragenbrink/slate/domain"

	"github.com/bwmarrin/discordgo"

	"github.com/stretchr/testify/suite"
)

type BotServiceTestSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	sess   *bot.MockDiscordSession
	now    time.Time
	router *bot.Router
}

func (s *BotServiceTestSuite) BeforeTest(suiteName, testName string) {
	s.ctrl = gomock.NewController(s.T())
	s.sess = bot.NewMockDiscordSession(s.ctrl)
	s.now = time.Now()
	cfg := &domain.SlateConfig{
		HerokuConfig: domain.HerokuConfig{
			ReleaseVersion:   "test",
			ReleaseCreatedAt: s.now.String(),
		},
	}
	s.router = SetupCommands(cfg)
	s.router.SetPrefix("$")
}

func (s *BotServiceTestSuite) AfterTest(suiteName, testName string) {
	s.ctrl.Finish()
}

func (s *BotServiceTestSuite) TestVersion() {
	msg := newMessageCreate("test", "$version")
	expected := fmt.Sprintf("version=test releaseDate=%s", s.now.String())
	wg := &sync.WaitGroup{}
	s.sess.EXPECT().ChannelMessageSend("test", expected).Do(func(a, b string) { wg.Done() })
	wg.Add(1)
	s.router.RouteMessageCreate(s.sess, msg)
	wg.Wait()
}

func newMessageCreate(user string, message string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID:       user,
				Username: user,
			},
			ChannelID: "test",
			Content:   message,
		},
	}
}

func TestBotService(t *testing.T) {
	suite.Run(t, new(BotServiceTestSuite))
}

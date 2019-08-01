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

package infrastructure

import (
	"fmt"
	"sync"

	"github.com/kkragenbrink/slate/interfaces"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/kkragenbrink/slate/domain"
)

// The DiscordServer controls interactions between our BotService router and our Discord session
type DiscordServer struct {
	cfg     *domain.SlateConfig
	discord *discordgo.Session
	run     *sync.Mutex
}

// NewDiscordServer initializes the bot service and discord session, then returns a DiscordServer object
func NewDiscordServer(cfg *domain.SlateConfig) *DiscordServer {
	connstr := fmt.Sprintf("Bot %s", cfg.DiscordToken)
	session, err := discordgo.New(connstr)

	router := interfaces.SetupCommands(cfg)
	router.SetPrefix(cfg.DiscordPrefix)
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) { router.RouteMessageCreate(s, m) })

	if err != nil {
		logrus.WithError(err).Fatal("could not establish discord session")
	}

	return &DiscordServer{
		cfg:     cfg,
		discord: session,
		run:     &sync.Mutex{},
	}
}

// Start locks the mutex and then connects to Discord in a goroutine
func (s *DiscordServer) Start() {
	s.run.Lock()
	go func() {
		if err := s.discord.Open(); err != nil {
			logrus.WithError(err).Fatal("could not connect to discord")
		}
	}()
}

// Stop closes the connection to Discord and unlocks the mutex
func (s *DiscordServer) Stop() {
	if err := s.discord.Close(); err != nil {
		logrus.WithError(err).Error("could not gracefully close connection to discord")
	}
	s.run.Unlock()
}

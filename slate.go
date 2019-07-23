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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kkragenbrink/slate/domain"
	"github.com/kkragenbrink/slate/interfaces"
)

// Slate is a discord bot and web API app
type Slate struct {
	bot    *interfaces.Bot
	config *domain.SlateConfig
	random *interfaces.Random
}

// NewSlate creates a new instance of the slate app
func NewSlate(config *domain.SlateConfig) (*Slate, error) {
	slate := &Slate{config: config}
	var err error

	slate.random = interfaces.NewRandom(slate.config)
	slate.bot, err = interfaces.NewBot(slate.config, slate.random)
	if err != nil {
		return nil, err
	}

	return slate, nil
}

// Start starts the services
func (s *Slate) Start() {
	s.bot.Start()
	waitForSignals()
}

// Stop stops the services
func (s *Slate) Stop() {
	s.bot.Stop()
}

func waitForSignals() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

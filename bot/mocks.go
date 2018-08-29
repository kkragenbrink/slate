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
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/config"
	"github.com/stretchr/testify/mock"
)

// MockDiscordSession mocks a discord session object
type MockDiscordSession struct {
	mock.Mock
}

// AddHandler is a stub
func (m *MockDiscordSession) AddHandler(h interface{}) func() {
	m.Called()
	return func() {}
}

// ChannelMessageSend is a stub
func (m *MockDiscordSession) ChannelMessageSend(ch string, message string) (*discordgo.Message, error) {
	m.Called()
	return nil, nil
}

// Open is a stub
func (m *MockDiscordSession) Open() error {
	m.Called()
	return nil
}

// Close is a stub
func (m *MockDiscordSession) Close() error {
	m.Called()
	return nil
}

// MockDiscordFactory is a DiscordFactory mock function for the MockDiscordSession
func MockDiscordFactory(cfg *config.Config) (DiscordSession, error) {
	return new(MockDiscordSession), nil
}

// MockCommand is a mock command handler
type MockCommand struct {
	mock.Mock
}

// Name returns the name of the command
func (*MockCommand) Name() string { return "mock" }

// Synopsis returns the synospis of the command
func (*MockCommand) Synopsis() string { return "used for test mocks" }

// Usage returns the usage of the command
func (*MockCommand) Usage() string { return "mock usage" }

// Execute runs the command
func (mc *MockCommand) Execute(ctx context.Context, fields []string, session DiscordSession, message *discordgo.MessageCreate) {
	mc.Called()
}

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

package config

import (
	"os"

	"github.com/pkg/errors"
)

// ErrNoDiscordToken is thrown when the environment variable isn't set
var ErrNoDiscordToken = errors.New("$DISCORD_TOKEN is required")

// Config holds configuration information which is necessary to run slate.
// This information is passed in at runtime via environment variables.
type Config struct {
	DiscordToken string
}

// New returns a new configuration object, and initializes th at object from
// environment variables.
func New() (*Config, error) {
	cfg := new(Config)

	// Initialize the Discord Token
	dt, err := initDiscordToken()
	if err != nil {
		return nil, err
	}
	cfg.DiscordToken = dt

	return cfg, nil
}

func initDiscordToken() (string, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return "", ErrNoDiscordToken
	}
	return token, nil
}

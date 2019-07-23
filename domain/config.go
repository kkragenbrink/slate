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

package domain

import (
	"github.com/caarlos0/env"
	"time"
)

// OAuthConfig holds the configuration information relevant to OAuth authentication
type OAuthConfig struct {
	ClientID     string `env:"OAUTH_CLIENT_ID" usage:"The client ID used when authenticating with OAUTH"`
	ClientSecret string `env:"OAUTH_CLIENT_SECRET" usage:"The client secret used when authenticating with OAUTH"`
}

// PostgresConfig holds configuration information relevant to the database connection
type PostgresConfig struct {
	Host     string `env:"DATABASE_HOST" envDefault:"localhost" usage:"The postgres host to connect to"`
	Port     int    `env:"DATABASE_PORT" envDefault:"5432" usage:"The postgres port to connect to"`
	Username string `env:"DATABASE_USER" usage:"The postgres username to authenticate with"`
	Password string `env:"DATABASE_PASSWORD" usage:"The postgres password to authenticate with"`
	Store    string `env:"DATABASE_STORE" usage:"The postgres store to query"`
}

// SlateConfig holds configuration information necessary to run slate
type SlateConfig struct {
	OAuthConfig
	PostgresConfig

	ApplicationHostName string        `env:"APPLICATION_HOSTNAME,required" usage:"The name of the application host, used in Authentication"`
	BotCommandPrefix    string        `env:"BOT_COMMAND_PREFIX" envDefault:"$" usage:"The prefix for commands recognized by the Discord bot"`
	BotCommandTimeout   time.Duration `env:"BOT_COMMAND_TIMEOUT" envDefault:"1s" usage:"The duration a command to the bot is allowed to wait before timing out"`
	Debug               bool          `env:"DEBUG" usage:"sets logrus loglevel to debug"`
	DiscordToken        string        `env:"DISCORD_TOKEN,required" usage:"The token to use when authenticating with Discord"`
	NodeID              int           `env:"NODE_ID" envDefault:"0" usage:"The ID of the current node (used for generating snowflake IDs)"`
	Port                int           `env:"PORT" envDefault:"3000" usage:"The port on which Slate listens for http requests"`
	RandomOrgAPIKey     string        `env:"RANDOM_ORG_API_KEY" usage:"The API key for requests to Random.org"`
	SessionSecret       string        `env:"SESSION_SECRET,required" usage:"Å randomly generated string for encrypting session information"`
}

func NewSlateConfig() (*SlateConfig, error) {
	config := &SlateConfig{}
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}

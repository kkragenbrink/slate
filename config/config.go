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
	"strconv"
)

// ErrNoDiscordToken is thrown when the environment variable isn't set
var ErrNoDiscordToken = errors.New("$DISCORD_TOKEN is required")

// ErrDatabaseInfo is thrown when the environment variables for the database aren't set
var ErrDatabaseInfo = errors.New("$DATABASE_HOST, $DATABASE_PORT, $DATABASE_USER, $DATABASE_PASS, and $DATABASE_NAME are required")

// Database holds configuration information for the database connection
type Database struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

// Config holds configuration information which is necessary to run slate.
// This information is passed in at runtime via environment variables.
type Config struct {
	CommandPrefix string
	Database      *Database
	DiscordToken  string
}

// New returns a new configuration object, and initializes th at object from
// environment variables.
func New() (*Config, error) {
	// Initialize the Discord Token
	discordToken, err := initDiscordToken()
	if err != nil {
		return nil, err
	}

	// Initialize the database
	database, err := initDatabase()
	if err != nil {
		return nil, err
	}

	// Initialize the Command Prefix
	commandPrefix := initCommandPrefix()

	cfg := &Config{
		DiscordToken:  discordToken,
		Database:      database,
		CommandPrefix: commandPrefix,
	}

	return cfg, nil
}

func initCommandPrefix() string {
	prefix := os.Getenv("COMMAND_PREFIX")
	if prefix == "" {
		return "$"
	}

	return prefix
}

func initDatabase() (*Database, error) {
	host := os.Getenv("DATABASE_HOST")
	portStr := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASS")
	name := os.Getenv("DATABASE_NAME")

	if host == "" || portStr == "" || user == "" || pass == "" || name == "" {
		return nil, ErrDatabaseInfo
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, ErrDatabaseInfo
	}

	return &Database{host, port, user, pass, name}, nil
}

func initDiscordToken() (string, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return "", ErrNoDiscordToken
	}
	return token, nil
}

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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew tests the New() function with all environment variables already set.
// This should return a valid Config{} struct and no errors.
func TestNew(t *testing.T) {
	// setup
	dt := os.Getenv("DISCORD_TOKEN")
	expectedDiscordToken := "test"
	os.Setenv("DISCORD_TOKEN", expectedDiscordToken)

	port := os.Getenv("PORT")
	expectedPort := 1234
	os.Setenv("PORT", strconv.Itoa(expectedPort))

	// run tests
	cfg, err := New()
	assert.Nil(t, err)
	assert.Equal(t, expectedDiscordToken, cfg.DiscordToken)
	assert.Equal(t, expectedPort, cfg.Port)

	// tear down
	os.Setenv("DISCORD_TOKEN", dt)
	os.Setenv("PORT", port)
}

// TestNoDiscordToken tests initDiscordToken without the environment variable.
func TestNoDiscordToken(t *testing.T) {
	dt, err := initDiscordToken()
	assert.Equal(t, "", dt)
	assert.Error(t, err)
}

// TestNoDiscordToken tests initDiscordToken without the environment variable.
func TestNoPort(t *testing.T) {
	port, err := initPort()
	assert.Equal(t, 0, port)
	assert.Error(t, err)
}

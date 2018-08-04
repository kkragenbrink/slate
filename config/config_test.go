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
	"testing"
)

// TestNew tests the New() function with all environment variables already set.
// This should return a valid Config{} struct and no errors.
func TestNew(t *testing.T) {
	// setup
	dt := os.Getenv("DISCORD_TOKEN")
	expectedDiscordToken := "test"
	os.Setenv("DISCORD_TOKEN", expectedDiscordToken)

	// run tests
	cfg, err := New()
	if err != nil {
		t.Errorf("New() returned error: %+v", err)
	}
	if cfg.DiscordToken != expectedDiscordToken {
		t.Errorf("New() did not correctly retrieve the correct environment variable")
	}

	// tear down
	os.Setenv("DISCORD_TOKEN", dt)
}

// TestNoDiscordToken tests initDiscordToken without the environment variable.
func TestNoDiscordToken(t *testing.T) {
	dt, err := initDiscordToken()
	if dt != "" {
		t.Error("initDiscordToken(), dt is set, should be nil")
	}
	if err == nil {
		t.Error("initDiscordToken(), err is nil, should be noDiscordToken error")
	}
}

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

package settings

import (
	"os"
	"testing"

	"strconv"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// setup
	dt := os.Getenv("DISCORD_TOKEN")
	expectedDiscordToken := "test"
	os.Setenv("DISCORD_TOKEN", expectedDiscordToken)
	dbhost := os.Getenv("DATABASE_HOST")
	dbport := os.Getenv("DATABASE_PORT")
	dbuser := os.Getenv("DATABASE_USER")
	dbpass := os.Getenv("DATABASE_PASS")
	dbname := os.Getenv("DATABASE_NAME")
	expectedDatabase := &Database{"localhost", 1234, "test", "test", "test"}
	os.Setenv("DATABASE_HOST", expectedDatabase.Host)
	os.Setenv("DATABASE_PORT", strconv.Itoa(expectedDatabase.Port))
	os.Setenv("DATABASE_USER", expectedDatabase.User)
	os.Setenv("DATABASE_PASS", expectedDatabase.Pass)
	os.Setenv("DATABASE_NAME", expectedDatabase.Name)
	oci := os.Getenv("OAUTH_CLIENT_ID")
	ocs := os.Getenv("OAUTH_CLIENT_SECRET")
	expectedOAuth := &OAuth{"test", "test"}
	os.Setenv("OAUTH_CLIENT_ID", expectedOAuth.ClientID)
	os.Setenv("OAUTH_CLIENT_SECRET", expectedOAuth.ClientSecret)
	ss := os.Getenv("SESSION_SECRET")
	expectedSessionSecret := "test"
	os.Setenv("SESSION_SECRET", expectedSessionSecret)

	// run tests
	cfg, err := Init()
	assert.Nil(t, err)
	assert.Equal(t, expectedDiscordToken, cfg.DiscordToken)
	assert.Equal(t, expectedDatabase, cfg.Database)

	// tear down
	os.Setenv("DISCORD_TOKEN", dt)
	os.Setenv("DATABASE_HOST", dbhost)
	os.Setenv("DATABASE_PORT", dbport)
	os.Setenv("DATABASE_USER", dbuser)
	os.Setenv("DATABASE_PASS", dbpass)
	os.Setenv("DATABASE_NAME", dbname)
	os.Setenv("OAUTH_CLIENT_ID", oci)
	os.Setenv("OAUTH_CLIENT_SECRET", ocs)
	os.Setenv("SESSION_SECRET", ss)
}

func TestCommandPrefix_Custom(t *testing.T) {
	// setup
	cp := os.Getenv("COMMAND_PREFIX")
	expectedCommandPrefix := "%"
	os.Setenv("COMMAND_PREFIX", expectedCommandPrefix)

	// run tests
	prefix := initCommandPrefix()
	assert.Equal(t, expectedCommandPrefix, prefix)

	// tear down
	os.Setenv("COMMAND_PREFIX", cp)
}

// TestDatabase_Missing test initDatabase without the required parameters.
func TestDatabase_Missing(t *testing.T) {
	// setup
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASS")
	name := os.Getenv("DATABASE_NAME")

	os.Setenv("DATABASE_HOST", "")
	os.Setenv("DATABASE_PORT", "")
	os.Setenv("DATABASE_USER", "")
	os.Setenv("DATABASE_PASS", "")
	os.Setenv("DATABASE_NAME", "")

	// run tests
	dc, err := initDatabase()
	assert.Error(t, err)

	os.Setenv("DATABASE_HOST", "localhost")
	dc, err = initDatabase()
	assert.Error(t, err)

	os.Setenv("DATABASE_PORT", "1234")
	dc, err = initDatabase()
	assert.Error(t, err)

	os.Setenv("DATABASE_USER", "test")
	dc, err = initDatabase()
	assert.Error(t, err)

	os.Setenv("DATABASE_PASS", "test")
	dc, err = initDatabase()
	assert.Error(t, err)

	os.Setenv("DATABASE_NAME", "test")
	expected := &Database{
		Host: "localhost",
		Port: 1234,
		User: "test",
		Pass: "test",
		Name: "test",
	}
	dc, err = initDatabase()
	assert.Equal(t, expected, dc)
	assert.Nil(t, err)

	// teardown
	os.Setenv("DATABASE_HOST", host)
	os.Setenv("DATABASE_PORT", port)
	os.Setenv("DATABASE_USER", user)
	os.Setenv("DATABASE_PASS", pass)
	os.Setenv("DATABASE_NAME", name)
}

func TestDiscordToken_Missing(t *testing.T) {
	dt, err := initDiscordToken()
	assert.Equal(t, "", dt)
	assert.Error(t, err)
}

func TestNodeID_Error(t *testing.T) {
	// setup
	envid := os.Getenv("NODE_ID")
	expectedNodeID := "five"
	os.Setenv("NODE_ID", expectedNodeID)

	// run tests
	id, err := initNodeID()
	assert.Equal(t, 0, id)
	assert.Error(t, err)

	// teardown
	os.Setenv("NODE_ID", envid)
}

func TestPort_Custom(t *testing.T) {
	// setup
	envport := os.Getenv("PORT")
	expectedPort := 1234
	os.Setenv("PORT", strconv.Itoa(expectedPort))

	// run tests
	port, err := initPort()
	assert.Equal(t, expectedPort, port)
	assert.Nil(t, err)

	// teardown
	os.Setenv("PORT", envport)
}

func TestPort_Error(t *testing.T) {
	// setup
	envport := os.Getenv("PORT")
	expectedPort := "test"
	os.Setenv("PORT", expectedPort)

	// run tests
	port, err := initPort()
	assert.Equal(t, 0, port)
	assert.Error(t, err)

	// teardown
	os.Setenv("PORT", envport)
}

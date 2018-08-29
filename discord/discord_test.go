package discord

import (
	"github.com/kkragenbrink/slate/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := new(config.Config)
	cfg.DiscordToken = "test-token"
	_, err := New(cfg)
	assert.Nil(t, err)
}

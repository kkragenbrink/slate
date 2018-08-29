package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/config"
)

// New returns a new session from discordgo, which matches the
// interface required by our bot package.
func New(cfg *config.Config) (*discordgo.Session, error) {
	connstr := fmt.Sprintf("Bot %s", cfg.DiscordToken)
	return discordgo.New(connstr)
}

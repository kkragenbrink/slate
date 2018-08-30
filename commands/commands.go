package commands

import (
	"github.com/kkragenbrink/slate/bot"
	"github.com/kkragenbrink/slate/commands/roll"
)

// Setup adds all available commands to the bot
func Setup(b *bot.Bot) {
	b.AddCommand(&roll.Command{})
}

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

package interfaces

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/domain"
	"github.com/kkragenbrink/slate/infrastructure/bot"
)

// The BotService contains handlers for bot commands
type BotService struct {
	Config *domain.SlateConfig
	Router *bot.Router
}

// SetupCommands creates a bot router for use in handling discord requests
func SetupCommands(cfg *domain.SlateConfig) *bot.Router {
	bs := &BotService{
		Config: cfg,
		Router: bot.NewRouter(),
	}

	bs.Router.AddCommand("version", bs.Version)

	return bs.Router
}

// Version handles version commands from discord
func (bs *BotService) Version(m *discordgo.Message) (response string) {
	defer formatResponse(m, response)

	response = fmt.Sprintf("version=%s releaseDate=%s",
		bs.Config.HerokuConfig.ReleaseVersion,
		bs.Config.HerokuConfig.ReleaseCreatedAt,
	)

	return
}

func formatResponse(m *discordgo.Message, response string) string {
	return fmt.Sprintf("%s %s", m.Author.Mention(), response)
}

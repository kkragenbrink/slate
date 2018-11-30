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

package sheet

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/kkragenbrink/slate/domains"
	"github.com/pkg/errors"
)

// New creates a new character sheet and stores it
func New(ctx context.Context, db domains.CharacterRepository, name, system string, guild, player int64) (*domains.Character, error) {
	character := new(domains.Character)
	character.Name = name
	character.Guild = guild
	character.Player = player
	character.System = system
	character.Sheet = GenerateSheetBySystem(system, nil)
	err := db.Store(ctx, character)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a new sheet")
	}
	return character, nil
}

// Get gets a character sheet from the database by ID
func Get(ctx context.Context, db domains.CharacterRepository, id snowflake.ID) (*domains.Character, error) {
	return db.FindByID(ctx, id.Int64())
}

// GenerateSheetBySystem generates a sheet by a specified system.  If the body is specified,
// this function will also populate that sheet from json
func GenerateSheetBySystem(system string, body json.RawMessage) domains.Sheet {
	switch {
	case system == "wtf2e":
		sheet := NewWtF2e()
		if body != nil {
			json.Unmarshal(body, sheet)
		}
		return sheet
	}
	return nil
}

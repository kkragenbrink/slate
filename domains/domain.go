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

package domains

import (
	"context"
	"github.com/bwmarrin/snowflake"
)

// A Character represents a character owned by a player
type Character struct {
	ID         *snowflake.ID `json:"id"`
	Name       string        `json:"name"`
	Guild      string        `json:"guild"`
	Player     string        `json:"player"`
	PlayerName string        `json:"playerName"`
	System     string        `json:"system"`
	Sheet      Sheet         `json:"sheet"`
}

// The CharacterRepository describes the interface to find and store characters.
type CharacterRepository interface {
	FindByPlayer(ctx context.Context, id string) ([]*Character, error)
	FindByID(ctx context.Context, id string) (*Character, error)
	Store(ctx context.Context, c *Character) error
}

// A Sheet is a type of character sheet.
type Sheet interface {
	System() string
}

// A User is an authenticated discord account.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

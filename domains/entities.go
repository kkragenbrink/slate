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

import "github.com/bwmarrin/snowflake"

// A Game is a collection of character sheets typically associated with a Discord Guild or Channel.
type Game struct {
	ID     snowflake.ID
	Name   string
	Player Player
}

// A GameRepository implements methods to store or find games by ID.
type GameRepository interface {
	// Persists a Game object.
	Store(game Game)
	// Finds a specific Game from the persistence storage.
	FindById(id snowflake.ID) Game
}

// A Player is a user who has owns one or more character sheets.
type Player struct {
	ID snowflake.ID
}

// A PlayerRepository implements methods to store or find players by ID.
type PlayerRepository interface {
	// Persists a Player object.
	Store(player Player)
	// Finds a specific Player from the persistence storage.
	FindByID(id snowflake.ID) Player
}

// A Sheet is a collection of object statistics which are associated with a game,
// and owned by a specific player.
type Sheet struct {
	ID     snowflake.ID
	Name   string
	Game   Game
	Player Player
}

// A SheetRepository implements methods to store or find sheets by ID.
type SheetRepository interface {
	// Persists a Sheet object.
	Store(sheet Sheet)
	// Finds a specific Sheet from the persistence storage.
	FindById(id snowflake.ID) Sheet
}

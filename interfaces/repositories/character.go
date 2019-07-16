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

package repositories

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/kkragenbrink/slate/domains"
	"github.com/kkragenbrink/slate/usecases/sheet"
	"github.com/pkg/errors"
	"strconv"
)

// The CharacterRepository stores the instructions to get and set from the database
type CharacterRepository struct {
	db Database
}

// NewCharacterRepository returns a new CharacterRepository instance
func NewCharacterRepository(db Database) *CharacterRepository {
	cr := new(CharacterRepository)
	cr.db = db
	return cr
}

// FindByPlayer retrieves a list of Characters from the database by the player ID.
func (cr *CharacterRepository) FindByPlayer(ctx context.Context, id string) ([]*domains.Character, error) {
	query := "SELECT id, name, guild, player, system FROM characters WHERE player = $1"
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse id")
	}
	rows, err := cr.db.Conn().QueryContext(ctx, query, pid)
	if err != nil {
		return nil, errors.Wrap(err, "could not get characters")
	}
	defer rows.Close()
	chars := make([]*domains.Character, 0)
	for rows.Next() {
		var char domains.Character
		var id int64
		err := rows.Scan(&id, &char.Name, &char.Guild, &char.Player, &char.System)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan character")
		}
		sid := snowflake.ID(id)
		char.ID = &sid
		chars = append(chars, &char)
	}
	return chars, nil
}

// FindByID retrieves a Character from the database by ID.
func (cr *CharacterRepository) FindByID(ctx context.Context, id string) (*domains.Character, error) {
	var c domains.Character
	var sid snowflake.ID
	var idc, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse id")
	}
	sid = snowflake.ID(idc)
	c.ID = &sid
	query := "SELECT name, guild, player, system, sheet FROM characters WHERE id = $1"
	row := cr.db.Conn().QueryRowContext(ctx, query, sid.Int64())
	var sh json.RawMessage
	err = row.Scan(&c.Name, &c.Guild, &c.Player, &c.System, &sh)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve character from the database")
	}
	c.Sheet = sheet.GenerateSheetBySystem(c.System, sh)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal sheet for character")
	}
	return &c, nil
}

// Store saves a character to the database.
// If the character does not yet have an ID (e.g. if it is new) it will create one at this point.
func (cr *CharacterRepository) Store(ctx context.Context, c *domains.Character) error {
	if c.ID == nil {
		c.ID = cr.db.ID()
	}
	sh, err := json.Marshal(c.Sheet)
	if err != nil {
		return errors.Wrap(err, "could not marshal sheet")
	}
	query := "INSERT INTO characters (id, name, guild, player, system, sheet) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, sheet = EXCLUDED.sheet"
	result, err := cr.db.Conn().ExecContext(ctx, query, c.ID.Int64(), c.Name, c.Guild, c.Player, c.System, sh)
	if err != nil {
		return errors.Wrap(err, "could not upsert character")
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		return errors.Wrap(err, "could not upsert character")
	}
	return nil
}

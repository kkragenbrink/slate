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

package database

import "github.com/jmoiron/sqlx/types"

// Metadata is used for storing program-level metadata, such as the database
// version.
type Metadata struct {
	Key   string         `db:"key"`
	Value types.JSONText `db:"value"`
}

var metadataSchema = `
CREATE TABLE IF NOT EXISTS metadata (
	key 	varchar(255)	NOT NULL 	PRIMARY KEY,
	value 	json 			NOT NULL
)`

// Sheet is used to store character sheets. It is keyed by a snowflake ID,
// but also tracks the UserID and GuildID of the creator.
type Sheet struct {
	ID      string         `db:"id"`
	UserID  string         `db:"userId"`
	GuildID string         `db:"guildId"`
	Name    string         `db:"name"`
	Sheet   types.JSONText `db:"sheet"`
}

var sheetSchema = `
CREATE TABLE IF NOT EXISTS sheet (
	ID 		bigint 			NOT NULL 	PRIMARY KEY,
	UserID 	bigint 			NOT NULL,
	GuildID bigint 			NOT NULL,
	Name 	varchar(255)	NOT NULL,
	Sheet 	json 			NOT NULL
)`

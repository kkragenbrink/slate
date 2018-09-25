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

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

func upgrade(db *sqlx.DB, version int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for v := version; v < DatabaseVersion; v++ {
		switch v {
		case 0:
			firstInstall(tx)
		}
	}
	return tx.Commit()
}

func firstInstall(tx *sql.Tx) {
	// Create tables
	tx.Exec(metadataSchema)
	tx.Exec(sheetSchema)

	// Set metadata
	version, _ := types.JSONText("1").Value()
	tx.Exec("INSERT INTO metadata (key, value) VALUES ($1, $2)", "DatabaseVersion", version)
}

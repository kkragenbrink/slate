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
	"fmt"
	"github.com/kkragenbrink/slate/config"

	"database/sql"
	"github.com/jmoiron/sqlx"
	// load drivers, but no need for pq itself
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"strconv"
)

// ErrNewDatabase is thrown when slate detects that it has connected to a new database.
var ErrNewDatabase = errors.New("new database discovered")

// ErrOutOfDate is thrown when the software version is behind the database version.
var ErrOutOfDate = errors.New("slate is out of date")

// DatabaseVersion is used internally to determine whether a database matches
// the current version of the application. Upgrades will be performed by
// iterating from one version to the next until it reaches the correct version.
var DatabaseVersion = 1

// New returns a new database connection, which matches the interface used by
// our sheet models.
func New(cfg *config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	err = validateVersion(db)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func validateVersion(db *sqlx.DB) error {
	var err error
	// Get the version from the database
	var vmeta Metadata
	err = db.Get(&vmeta, "SELECT value FROM metadata WHERE key=$1", "DatabaseVersion")
	if err != nil && err != sql.ErrNoRows {
		if err.Error() == `pq: relation "metadata" does not exist` {
			// new database
			return upgrade(db, 0)
		}
		return err
	}

	var version int
	if vmeta.Value != nil {
		version, err = strconv.Atoi(vmeta.Value.String())
		if err != nil {
			return err
		}
	}

	// Compare to the software version expectations
	if version < DatabaseVersion {
		return upgrade(db, version)
	}

	if version > DatabaseVersion {
		return ErrOutOfDate
	}
	return nil
}

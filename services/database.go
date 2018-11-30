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

package services

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/kkragenbrink/slate/interfaces/repositories"
	"github.com/kkragenbrink/slate/settings"
	// postgres bindings
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// The DatabaseService manages the lifecycle of the connection to postgres
type DatabaseService struct {
	conn     *sql.DB
	flake    *snowflake.Node
	settings *settings.Settings
	repos    map[string]interface{}
}

// NewDatabaseService creates a new instance of the DatabaseService
func NewDatabaseService(set *settings.Settings) *DatabaseService {
	dbs := new(DatabaseService)
	dbs.settings = set
	dbs.initModels()
	return dbs
}

// Conn returns the actual database connection for sql methods
func (dbs *DatabaseService) Conn() *sql.DB {
	return dbs.conn
}

func (dbs *DatabaseService) initModels() {
	dbs.repos = make(map[string]interface{})
	dbs.repos["character"] = repositories.NewCharacterRepository(dbs)
}

// Repository retrieves a specific repository by name
func (dbs *DatabaseService) Repository(name string) interface{} {
	return dbs.repos[name]
}

// ID generates a new primary key ID using snowflake
func (dbs *DatabaseService) ID() *snowflake.ID {
	id := dbs.flake.Generate()
	return &id
}

// Start establishes a connection to Postgres
func (dbs *DatabaseService) Start() error {
	connstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		dbs.settings.Database.Host, dbs.settings.Database.Port, dbs.settings.Database.User,
		dbs.settings.Database.Pass, dbs.settings.Database.Name)
	var err error
	dbs.conn, err = sql.Open("postgres", connstr)
	dbs.flake, err = snowflake.NewNode(int64(dbs.settings.NodeID))
	if err != nil {
		return errors.Wrap(err, "could not connect to postgres")
	}
	return nil
}

// Stop closes the connection to Postgres
func (dbs *DatabaseService) Stop() error {
	return dbs.conn.Close()
}

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
	"github.com/kkragenbrink/slate/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DatabaseSuite struct {
	suite.Suite
}

func TestDatabaseSuite(t *testing.T) {
	s := new(DatabaseSuite)
	suite.Run(t, s)
}

func newMockDBSettings() *settings.Settings {
	set := new(settings.Settings)
	database := new(settings.Database)
	database.User = "test"
	database.Name = "test"
	database.Host = "test"
	database.Pass = "test"
	database.Port = 1234
	set.Database = database
	return set
}

func (suite *DatabaseSuite) TestNewDatabaseService() {
	db := NewDatabaseService(newMockDBSettings())
	assert.NotNil(suite.T(), db)
}

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

package main

import (
	"fmt"
	"os"

	"github.com/kkragenbrink/slate/services"
	"github.com/kkragenbrink/slate/settings"
)

func main() {
	// Initialize the settings
	set, err := settings.Init()
	handleError(err, 1)

	// Create Services
	db := services.NewDatabaseService(set)
	rand := services.NewRandom(set)
	auth := services.NewAuthService(set)
	bot, err := services.NewBot(set, db, rand)
	handleError(err, 1)
	ws := services.NewWebService(set, auth, bot, db, rand)

	// Start services
	sm := NewServicesManager(db, bot, ws)
	sm.Start()
}

func handleError(err error, code int) {
	if err != nil {
		fmt.Print(err)
		os.Exit(code)
	}
}

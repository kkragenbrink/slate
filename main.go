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

package main

import (
	"fmt"
	"github.com/kkragenbrink/slate/services"
	"github.com/kkragenbrink/slate/settings"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize the settings
	set, err := settings.Init()
	handleError(err, 1)

	// Create Services
	bot, err := services.NewBot(set)
	handleError(err, 1)
	ws := services.NewWebService(set)

	// Start Services
	err = bot.Start()
	handleError(err, 1)
	err = ws.Start()
	handleError(err, 1)

	// Wait for signals
	waitForSignals()

	// Shutdown services
	err = ws.Stop()
	handleError(err, 2)
	err = bot.Stop()
	handleError(err, 2)
}

func handleError(err error, code int) {
	if err != nil {
		fmt.Print(err)
		os.Exit(code)
	}
}

func waitForSignals() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

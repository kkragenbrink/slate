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

package interfaces

import (
	"context"
	"encoding/json"
	"github.com/kkragenbrink/slate/usecases/roll"
)

// A Router contains methods for adding methods to handle http requests
type Router interface {
	HandlePost(string, WebHandler)
}

// A WebHandler is an http request handler
type WebHandler func(ctx context.Context, body json.RawMessage) (interface{}, error)

// A WebRoll is a roll command, and is used to determine the system
type WebRoll struct {
	System string `json:"system"`
}

// WebRollHandler handles incoming roll messages and sends them to the roll usecase.
func WebRollHandler(ctx context.Context, body json.RawMessage) (interface{}, error) {
	var r WebRoll
	err := json.Unmarshal(body, &r)
	rs, err := roll.NewRoller(r.System, body)
	if err != nil {
		// todo: log
		return 400, roll.ErrInvalidRollSystem
	}
	err = rs.Roll(ctx, nil)
	if err != nil {
		return 500, err
	}
	return rs, nil
}

func initWebServices(router Router) {
	router.HandlePost("/roll", WebRollHandler)
}

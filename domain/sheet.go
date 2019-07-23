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

package domain

// A SheetSystem represents the data structure for the character sheet
type SheetSystem string

// A Character represents a metadata storage object for a Sheet
type Character struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	GuildID   string      `json:"guildID"`
	ChannelID string      `json:"channelID"`
	PlayerID  string      `json:"playerID"`
	Player    string      `json:"player"`
	System    SheetSystem `json:"system"`
	Sheet     Sheet       `json:"sheet"`
}

// A Sheet represents a character sheet and all of its data
type Sheet interface {
	System() SheetSystem
}

// A SheetNote represents a note related to a character sheet
type SheetNote struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

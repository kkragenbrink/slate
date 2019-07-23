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

// A RollSystem is a series of rules regarding how dice should be interpreted
type RollSystem string

const (
	// RollSystemD20 represents the d20 dice system
	RollSystemD20 RollSystem = "d20"

	// RollSystemStoryteller represents the Storyteller (d10) dice system
	RollSystemStoryteller RollSystem = "storyteller"
)

// A Roll represents a request for the system to create a randomized outcome
type Roll struct {
	Expression string     `json:"expression"`
	System     RollSystem `json:"system"`
	Verbose    bool       `json:"verbose"`
}

// A D20Roll is a specific Roll for using the d20 dice system
type D20Roll struct {
	Roll
	Tokens []D20RollToken `json:"tokens"`
}

// A D20RollToken represents a tokenized portion of a D20Roll
type D20RollToken struct {
	Dice        int   `json:"dice"`
	Sides       int64 `json:"sides"`
	KeepHighest int   `json:"keepHighest"`
	KeepLowest  int   `json:"keepLowest"`
}

// A StorytellerRoll is a specific Roll for using the Storyteller dice system
type StorytellerRoll struct {
	Roll
	Again       int64 `json:"again"`
	Exceptional int   `json:"exceptional"`
	Rote        bool  `json:"rote"`
	Target      int   `json:"target"`
	Weakness    bool  `json:"weakness"`
}

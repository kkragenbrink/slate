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

package sheet

// CofD2eSpirit describes a CofD character sheet for ephemeral entities
type CofD2eSpirit struct {
	*BaseCofD2e

	// Core
	Rank   string `json:"rank"`
	Type   string `json:"type"`
	Vice   string `json:"vice"`
	Virtue string `json:"virtue"`

	// Attributes
	Power      int `json:"power"`
	Finesse    int `json:"finesse"`
	Resistance int `json:"resistance"`

	// Traits
	Essence        IntWithMax    `json:"essence"`
	Integrity      int           `json:"integrity"`
	Numina         []string      `json:"numina"`
	Anchors        []string      `json:"anchors"`
	Influences     []CofD2eMerit `json:"influences"`
	Manifestations []string      `json:"manifestations"`
	Ban            string        `json:"ban"`
	Bane           string        `json:"bane"`
}

// NewCofD2eSpirit creates a new instance of a CofD2e ephemeral sheet
func NewCofD2eSpirit() *CofD2eSpirit {
	sheet := new(CofD2eSpirit)
	sheet.BaseCofD2e = NewBaseCofD2e()
	sheet.Power = 1
	sheet.Finesse = 1
	sheet.Resistance = 1
	sheet.Numina = make([]string, 0)
	sheet.Anchors = make([]string, 0)
	sheet.Influences = make([]CofD2eMerit, 0)
	sheet.Manifestations = make([]string, 0)
	return sheet
}

// System returns the type of the sheet
func (s *CofD2eSpirit) System() string {
	return "cofd2e-spirit"
}

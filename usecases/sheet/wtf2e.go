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

// WtF2e describes a sheet for Werewolf the Forsaken
type WtF2e struct {
	*BaseCofD2e
	*CofD2eCreature

	// Core
	Blood   string `json:"blood"`
	Bone    string `json:"bone"`
	Auspice string `json:"auspice"`
	Tribe   string `json:"tribe"`
	Lodge   string `json:"lodge"`

	// Renown
	Cunning int `json:"cunning"`
	Glory   int `json:"glory"`
	Honor   int `json:"honor"`
	Purity  int `json:"purity"`
	Wisdom  int `json:"wisdom"`

	// Gifts and Rites
	Gifts struct {
		Moon   []WtF2eMoonGift `json:"moon"`
		Shadow []string        `json:"shadow"`
		Wolf   []string        `json:"wolf"`
	} `json:"gifts"`
	Rites []string `json:"rites"`

	// Traits
	Essence        IntWithMax `json:"essence"`
	Harmony        int        `json:"harmony"`
	KuruthTriggers string     `json:"kuruth_triggers"`
	PrimalUrge     int        `json:"primal_urge"`
	Touchstones    struct {
		Flesh  string `json:"flesh"`
		Spirit string `json:"spirit"`
	} `json:"touchstones"`
}

// WtF2eMoonGift describes a gift from the moon gift list
type WtF2eMoonGift struct {
	List string `json:"list"`
	Dots int    `json:"dots"`
}

// NewWtF2e creates a new instance of WtF2e and initializes its values
func NewWtF2e() *WtF2e {
	sheet := new(WtF2e)
	sheet.BaseCofD2e = NewBaseCofD2e()
	sheet.CofD2eCreature = NewCofD2eCreature()
	sheet.Gifts.Moon = []WtF2eMoonGift{
		{List: "", Dots: 0},
		{List: "", Dots: 0},
	}
	sheet.Gifts.Shadow = make([]string, 0)
	sheet.Gifts.Wolf = make([]string, 0)
	sheet.Rites = make([]string, 0)
	return sheet
}

// System identifies the sheet as a wtf2e sheet
func (s *WtF2e) System() string {
	return "wtf2e"
}

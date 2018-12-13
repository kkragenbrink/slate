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

// BaseCofD2e describes the base of every CofD2e sheet
type BaseCofD2e struct {
	// Core
	Name        string `json:"name"`
	Chronicle   string `json:"chronicle"`
	Concept     string `json:"concept"`
	Experiences int    `json:"experiences"`
	Beats       int    `json:"beats"`
	Notes       []Note `json:"notes"`

	// Traits
	Aspirations []string `json:"aspirations"`
	Conditions  string   `json:"conditions"`
	Health      struct {
		Max        int `json:"max"`
		Aggravated int `json:"aggravated"`
		Lethal     int `json:"lethal"`
		Bashing    int `json:"bashing"`
	} `json:"health"`
	Merits    []CofD2eMerit `json:"merits"`
	Size      int           `json:"size"`
	Willpower IntWithMax    `json:"willpower"`
}

// CofD2eCreature describes attributes and skills for a non-ephemeral entity
type CofD2eCreature struct {
	// Attributes
	Intelligence int `json:"intelligence"`
	Wits         int `json:"wits"`
	Resolve      int `json:"resolve"`
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Stamina      int `json:"stamina"`
	Presence     int `json:"presence"`
	Manipulation int `json:"manipulation"`
	Composure    int `json:"composure"`

	// Skills
	Academics     CofD2eSkill `json:"academics"`
	Computer      CofD2eSkill `json:"computer"`
	Crafts        CofD2eSkill `json:"crafts"`
	Investigation CofD2eSkill `json:"investigation"`
	Medicine      CofD2eSkill `json:"medicine"`
	Occult        CofD2eSkill `json:"occult"`
	Politics      CofD2eSkill `json:"politics"`
	Science       CofD2eSkill `json:"science"`
	Athletics     CofD2eSkill `json:"athletics"`
	Brawl         CofD2eSkill `json:"brawl"`
	Drive         CofD2eSkill `json:"drive"`
	Firearms      CofD2eSkill `json:"firearms"`
	Larceny       CofD2eSkill `json:"larceny"`
	Stealth       CofD2eSkill `json:"stealth"`
	Survival      CofD2eSkill `json:"survival"`
	Weaponry      CofD2eSkill `json:"weaponry"`
	AnimalKen     CofD2eSkill `json:"animal_ken"`
	Empathy       CofD2eSkill `json:"empathy"`
	Expression    CofD2eSkill `json:"expression"`
	Intimidation  CofD2eSkill `json:"intimidation"`
	Persuasion    CofD2eSkill `json:"persuasion"`
	Socialize     CofD2eSkill `json:"socialize"`
	Streetwise    CofD2eSkill `json:"streetwise"`
	Subterfuge    CofD2eSkill `json:"subterfuge"`
}

// CofD2e describes a sheet for a Mortal
type CofD2e struct {
	*BaseCofD2e
	*CofD2eCreature

	// Core
	Vice     string `json:"vice"`
	Virtue   string `json:"virtue"`
	Faction  string `json:"faction"`
	Group    string `json:"group"`
	Template string `json:"template"`

	// Traits
	Integrity int `json:"integrity"`
}

// CofD2eSkill describes a skill within the CofD 2e system
type CofD2eSkill struct {
	Dots        int    `json:"dots"`
	Specialties string `json:"specialties"`
}

// CofD2eMerit describes a merit within the CofD 2e system
type CofD2eMerit struct {
	Name string `json:"name"`
	Dots int    `json:"dots"`
}

// NewBaseCofD2e creates a new instance of a base CofD2e sheet, and initializes many of its values.
func NewBaseCofD2e() *BaseCofD2e {
	sheet := new(BaseCofD2e)
	sheet.Size = 5
	sheet.Notes = make([]Note, 0)
	sheet.Aspirations = make([]string, 0)
	sheet.Merits = make([]CofD2eMerit, 0)
	return sheet
}

// NewCofD2eCreature creates a new instance of a CofD2e creature sheet, and initializes many of its values.
func NewCofD2eCreature() *CofD2eCreature {
	sheet := new(CofD2eCreature)
	sheet.Intelligence = 1
	sheet.Wits = 1
	sheet.Resolve = 1
	sheet.Strength = 1
	sheet.Dexterity = 1
	sheet.Stamina = 1
	sheet.Presence = 1
	sheet.Manipulation = 1
	sheet.Composure = 1
	return sheet
}

// NewCofD2e creates a new instance of the mortal sheet and initializes its values
func NewCofD2e() *CofD2e {
	sheet := new(CofD2e)
	sheet.BaseCofD2e = NewBaseCofD2e()
	sheet.CofD2eCreature = NewCofD2eCreature()
	return sheet
}

// System returns the system of the sheet
func (s *CofD2e) System() string {
	return "cofd2e"
}

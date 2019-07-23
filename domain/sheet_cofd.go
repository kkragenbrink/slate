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

// A CofDSheet represents any character using the Chronicles of Darkness sheet system
type CofDSheet struct {
	Name             string            `json:"name"`
	Chronicle        string            `json:"chronicle"`
	Concept          string            `json:"concept"`
	Experiences      int               `json:"experiences"`
	SpentExperiences int               `json:"spent_experiences"`
	Beats            int               `json:"beats"`
	Aspirations      []string          `json:"apsirations"`
	Conditions       []string          `json:"conditions"`
	Size             CofDSize          `json:"size"`
	Willpower        CofDResourceTrack `json:"willpower"`
	Notes            []SheetNote       `json:"notes"`
}

// A CofDHealthTrack allows tracking the current and max health of a character
type CofDHealthTrack struct {
	Aggravated int            `json:"aggravated"`
	Lethal     int            `json:"lethal"`
	Bashing    int            `json:"bashing"`
	Modifiers  []CofDModifier `json:"modifiers"`
}

// A CofDModifier represents a modifier to a resource track or value
type CofDModifier struct {
	Value  int    `json:"value"`
	Source string `json:"source"`
}

// A CofDSize allows tracking the size of a character
type CofDSize struct {
	Base      int            `json:"base"`
	Modifiers []CofDModifier `json:"modifiers"`
}

// A CofDResourceTrack allows tracking the current resources for a character, and modifiers to the maximum value
type CofDResourceTrack struct {
	Current   int            `json:"current"`
	Modifiers []CofDModifier `json:"modifiers"`
}

// A CofDCorporeal describes attributes and skills which all corporeal entities have
type CofDCorporeal struct {
	// Attributes
	Intelligence CofDAttribute `json:"intelligence"`
	Wits         CofDAttribute `json:"wits"`
	Resolve      CofDAttribute `json:"resolve"`
	Strength     CofDAttribute `json:"strength"`
	Dexterity    CofDAttribute `json:"dexterity"`
	Stamina      CofDAttribute `json:"stamina"`
	Presence     CofDAttribute `json:"presence"`
	Manipulation CofDAttribute `json:"manipulation"`
	Composure    CofDAttribute `json:"composure"`

	// Skills
	Academics     CofDSkill `json:"academics"`
	Computer      CofDSkill `json:"computer"`
	Crafts        CofDSkill `json:"crafts"`
	Investigation CofDSkill `json:"investigation"`
	Medicine      CofDSkill `json:"medicine"`
	Occult        CofDSkill `json:"occult"`
	Politics      CofDSkill `json:"politics"`
	Science       CofDSkill `json:"science"`
	Athletics     CofDSkill `json:"athletics"`
	Brawl         CofDSkill `json:"brawl"`
	Drive         CofDSkill `json:"drive"`
	Firearms      CofDSkill `json:"firearms"`
	Larceny       CofDSkill `json:"larceny"`
	Stealth       CofDSkill `json:"stealth"`
	Survival      CofDSkill `json:"survival"`
	Weaponry      CofDSkill `json:"weaponry"`
	AnimalKen     CofDSkill `json:"animalKen"`
	Empathy       CofDSkill `json:"empathy"`
	Expression    CofDSkill `json:"expression"`
	Intimidation  CofDSkill `json:"intimidation"`
	Persuasion    CofDSkill `json:"persuasion"`
	Socialize     CofDSkill `json:"socialize"`
	Streetwise    CofDSkill `json:"streetwise"`
	Subterfuge    CofDSkill `json:"subterfuge"`

	// Other Stuff
	Merits []CofDMerit     `json:"merit"`
	Health CofDHealthTrack `json:"health"`
}

// A CofDEphemeral is any non-corporeal entity (Spirits, Ghosts, Angels, etc ...)
type CofDEphemeral struct {
	// Attributes
	Power      CofDAttribute `json:"power"`
	Finesse    CofDAttribute `json:"finesse"`
	Resistance CofDAttribute `json:"resistance"`

	// Other Stuff
	Corpus  CofDHealthTrack   `json:"corpus"`
	Essence CofDResourceTrack `json:"essence"`
}

// A CofDAttribute represents an attribute
type CofDAttribute struct {
	Dots      int            `json:"dots"`
	Modifiers []CofDModifier `json:"modifiers"`
}

// A CofDSkill represents a skill
type CofDSkill struct {
	Dots        int            `json:"dots"`
	Specialties []string       `json:"specialties"`
	Modifiers   []CofDModifier `json:"modifiers"`
}

// A CofDMerit represents a merit
type CofDMerit struct {
	Name string `json:"name"`
	Dots int    `json:"dots"`
}

////////////////////////////////////////
// Sheet Implementations
////////////////////////////////////////

// Constants representing the list of Chronicles of Darkness sheet types
const (
	SheetSystemChroniclesOfDarknessMortal   SheetSystem = "cofdMortal"
	SheetSystemChroniclesOfDarknessSpirit   SheetSystem = "cofdSpirit"
	SheetSystemChroniclesOfDarknessWerewolf SheetSystem = "cofdWerewolf"
)

// -- MORTAL --

// A CofDMortal is a Mortal or Mortal+
type CofDMortal struct {
	CofDSheet
	CofDCorporeal

	Vice      string `json:"vice"`
	Virtue    string `json:"virtue"`
	Faction   string `json:"faction"`
	Group     string `json:"group"`
	Integrity int    `json:"integrity"`

	ResourceTracks []struct {
		Name  string            `json:"name"`
		Track CofDResourceTrack `json:"track"`
	} `json:"resourceTracks"`
}

// System returns the CofD Mortal system
func (CofDMortal) System() SheetSystem {
	return SheetSystemChroniclesOfDarknessMortal
}

// -- WEREWOLF --

// A CofDWerewolf is a Werewolf
type CofDWerewolf struct {
	CofDSheet
	CofDCorporeal

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
	GiftLists struct {
		Moon   []CofDWerewolfMoonGiftList `json:"moon"`
		Shadow []CofDWerewolfGiftList     `json:"shadow"`
		Wolf   []CofDWerewolfGiftList     `json:"wolf"`
	} `json:"gift-lists"`

	// Traits
	Essence        CofDResourceTrack `json:"essence"`
	Harmony        int               `json:"harmony"`
	KuruthTriggers string            `json:"kuruthTriggers"`
	PrimalUrge     int               `json:"primalUrge"`
	TouchStones    struct {
		Flesh  []string `json:"flesh"`
		Spirit []string `json:"spirit"`
	} `json:"touchstones"`
}

// System returns the CofD Werewolf system
func (CofDWerewolf) System() SheetSystem {
	return SheetSystemChroniclesOfDarknessWerewolf
}

// A CofDWerewolfMoonGiftList represents a moon gift list in the Werewolf system
type CofDWerewolfMoonGiftList struct {
	CofDWerewolfGiftList
	Dots int `json:"dots"`
}

// A CofDWerewolfGiftList represents a gift list in the Werewolf system
type CofDWerewolfGiftList struct {
	Name  string             `json:"name"`
	Gifts []CofDWerewolfGift `json:"gifts"`
}

// A CofDWerewolfGift represents a gift in the Werewolf system
type CofDWerewolfGift struct {
	Name   string `json:"name"`
	Renown string `json:"renown"`
}

// -- SPIRITS --

// A CofDSpirit is a Spirit
type CofDSpirit struct {
	CofDSheet
	CofDEphemeral

	Rank           int         `json:"rank"`
	Choir          string      `json:"choir"`
	Court          string      `json:"court"`
	Ban            string      `json:"ban"`
	Bane           string      `json:"bane"`
	Numina         []string    `json:"numina"`
	Influences     []CofDMerit `json:"influences"`
	Manifestations []string    `json:"manifestations"`
}

// System returns the CofD Spirit system
func (CofDSpirit) System() SheetSystem {
	return SheetSystemChroniclesOfDarknessSpirit
}

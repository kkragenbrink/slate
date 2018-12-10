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
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi"
	"github.com/kkragenbrink/slate/domains"
	"github.com/kkragenbrink/slate/interfaces/repositories"
	"github.com/kkragenbrink/slate/usecases/roll"
	"github.com/kkragenbrink/slate/usecases/sheet"
	"github.com/kkragenbrink/slate/util"
	"io/ioutil"
	"net/http"
	"strconv"
)

// The WebServiceHandler stores information useful to the web service routes
type WebServiceHandler struct {
	bot Bot
	db  repositories.Database
}

// NewWebServiceHandler creates a new WebServiceHandler instance
func NewWebServiceHandler(bot Bot, db repositories.Database) *WebServiceHandler {
	ws := new(WebServiceHandler)
	ws.bot = bot
	ws.db = db
	return ws
}

// A Roll is a roll command, and is used to determine the system
type Roll struct {
	System string `json:"system"`
}

// Roll handles the roll usecase from the web.
func (ws *WebServiceHandler) Roll(res http.ResponseWriter, req *http.Request) {
	// decode
	body, err := util.Decodejson(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	req.Body.Close()
	// find the system
	var r Roll
	err = json.Unmarshal(body, &r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// get a roller
	rs, err := roll.NewRoller(r.System, body)
	if err != nil {
		http.Error(res, roll.ErrInvalidRollSystem.Error(), http.StatusBadRequest)
		return
	}
	// roll
	err = rs.Roll(req.Context(), nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// encode the results
	err = json.NewEncoder(res).Encode(rs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// A WebCharacter allows easy inspection of a character before unmarshalling the sheet.
type WebCharacter struct {
	Name  string          `json:"name"`
	Sheet json.RawMessage `json:"sheet"`
}

// Sheet handles the sheet usecase from the web.
func (ws *WebServiceHandler) Sheet(res http.ResponseWriter, req *http.Request) {
	rawid := chi.URLParam(req, "ID")
	parsed, err := strconv.ParseInt(rawid, 10, 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	id := snowflake.ID(parsed)
	repo := ws.db.Repository("character").(domains.CharacterRepository)
	char, err := sheet.Get(req.Context(), repo, id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	if req.Method == http.MethodPost {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		req.Body.Close()
		var tmp WebCharacter
		err = json.Unmarshal(body, &tmp)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		char.Name = tmp.Name
		sh := sheet.GenerateSheetBySystem(char.System, tmp.Sheet)
		char.Sheet = sh
		err = repo.Store(req.Context(), char)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
	if req.Method == http.MethodGet {
		user, err := ws.bot.User(strconv.FormatInt(char.Player, 10))
		if err != nil {
			// todo: log it
			char.PlayerName = err.Error()
		} else {
			char.PlayerName = user.String()
		}
	}
	err = json.NewEncoder(res).Encode(char)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

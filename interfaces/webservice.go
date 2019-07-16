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
	"fmt"
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
	auth Auth
	bot  Bot
	db   repositories.Database
	rand Random
}

// NewWebServiceHandler creates a new WebServiceHandler instance
func NewWebServiceHandler(auth Auth, bot Bot, db repositories.Database, rand Random) *WebServiceHandler {
	ws := new(WebServiceHandler)
	ws.auth = auth
	ws.bot = bot
	ws.db = db
	ws.rand = rand
	return ws
}

// A ChannelQuery is a channel command, and is used to determine the guild
type ChannelQuery struct {
	Guild    string     `json:"guild"`
	Channels []*Channel `json:"channels"`
}

// A Channel channel
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Channels lists the possible channels that can be rolled to.
func (ws *WebServiceHandler) Channels(res http.ResponseWriter, req *http.Request) {
	// decode
	defer req.Body.Close()
	body, err := util.Decodejson(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	var q ChannelQuery
	err = json.Unmarshal(body, &q)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	channels, err := ws.bot.Channels(q.Guild)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ch := range channels {
		channel := &Channel{
			ID:   ch.ID,
			Name: ch.Name,
		}
		q.Channels = append(q.Channels, channel)
	}
	err = json.NewEncoder(res).Encode(q)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// A Roll is a roll command, and is used to determine the system
type Roll struct {
	System  string `json:"system"`
	Channel string `json:"channel"`
}

// Roll handles the roll usecase from the web.
func (ws *WebServiceHandler) Roll(res http.ResponseWriter, req *http.Request) {
	if !ws.auth.IsAuthorized(req) {
		res.WriteHeader(http.StatusForbidden)
		return
	}
	user, err := ws.auth.GetAuthorization(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// decode
	defer req.Body.Close()
	body, err := util.Decodejson(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
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
	rs.SetRand(ws.rand.Rand)

	// roll
	err = rs.Roll(req.Context(), nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// send it to discord
	err = ws.bot.SendMessage(r.Channel, fmt.Sprintf("From the web: <@%s> %s", strconv.FormatInt(user.ID, 10), rs.ToString()))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
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

// Characters handles the request for all characters for a user from the web.
func (ws *WebServiceHandler) Characters(res http.ResponseWriter, req *http.Request) {
	if !ws.auth.IsAuthorized(req) {
		res.WriteHeader(http.StatusForbidden)
		return
	}
	user, err := ws.auth.GetAuthorization(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	repo := ws.db.Repository("character").(domains.CharacterRepository)
	chars, err := repo.FindByPlayer(req.Context(), strconv.FormatInt(user.ID, 10))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(res).Encode(chars)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// Sheet handles the sheet usecase from the web.
func (ws *WebServiceHandler) Sheet(res http.ResponseWriter, req *http.Request) {
	if !ws.auth.IsAuthorized(req) {
		res.WriteHeader(http.StatusForbidden)
		return
	}
	user, err := ws.auth.GetAuthorization(req)
	rawid := chi.URLParam(req, "ID")
	parsed, err := strconv.ParseInt(rawid, 10, 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	id := snowflake.ID(parsed)
	repo := ws.db.Repository("character").(domains.CharacterRepository)
	char, err := sheet.Get(req.Context(), repo, id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.Method == http.MethodPost {
		if strconv.FormatInt(user.ID, 10) != char.Player {
			http.Error(res, err.Error(), http.StatusForbidden)
			return
		}
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		var tmp WebCharacter
		err = json.Unmarshal(body, &tmp)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		char.Name = tmp.Name
		sh := sheet.GenerateSheetBySystem(char.System, tmp.Sheet)
		char.Sheet = sh
		err = repo.Store(req.Context(), char)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if req.Method == http.MethodGet {
		user, err := ws.bot.User(char.Player)
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

// Auth describes the interface for authorization of this application
type Auth interface {
	BeginAuthorization(http.ResponseWriter, *http.Request)
	CompleteAuthorization(http.ResponseWriter, *http.Request)
	GetAuthorization(*http.Request) (*domains.User, error)
	IsAuthorized(*http.Request) bool
}

// Auth displays the current authorization information about the user.
func (ws *WebServiceHandler) Auth(res http.ResponseWriter, req *http.Request) {
	if !ws.auth.IsAuthorized(req) {
		res.WriteHeader(http.StatusForbidden)
		return
	}
	user, err := ws.auth.GetAuthorization(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(res).Encode(user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// AuthBegin handles the usecase for attempting to log in from the web.
func (ws *WebServiceHandler) AuthBegin(res http.ResponseWriter, req *http.Request) {
	if ws.auth.IsAuthorized(req) {
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
	}

	ws.auth.BeginAuthorization(res, req)
}

// AuthComplete completes authorization and redirects the user back to the homepage.
func (ws *WebServiceHandler) AuthComplete(res http.ResponseWriter, req *http.Request) {
	if !ws.auth.IsAuthorized(req) {
		ws.auth.CompleteAuthorization(res, req)
	}

	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}

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

package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
	"github.com/kkragenbrink/slate/domains"
	"github.com/kkragenbrink/slate/interfaces/oauth2"
	"github.com/kkragenbrink/slate/settings"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	// The name of the session in the session store
	sessionName = "_slate_session"

	// The name of the state in the session store
	stateValue = "State"

	// The user object stored in the session store
	userValue = "User"

	// the identity scope for the discord api
	identityScope = "identify"
)

// ErrBadState is thrown when an invalid authorization attempt is made
var ErrBadState = errors.New("state does not match expected state")

// The AuthService handles requests to authorize a user
type AuthService struct {
	OAuthConfig *oauth2.Config
	Store       sessions.Store
}

// NewAuthService returns a new instance of the AuthService
func NewAuthService(set *settings.Settings) *AuthService {
	auth := new(AuthService)
	auth.initOAuth(set)
	auth.initStore(set)
	return auth
}

func (a *AuthService) initOAuth(set *settings.Settings) {
	a.OAuthConfig = &oauth2.Config{
		ClientID:     set.OAuth.ClientID,
		ClientSecret: set.OAuth.ClientSecret,
		Scopes:       []string{identityScope},
		Endpoint:     discord.Endpoint,
		RedirectURL:  fmt.Sprintf("https://%s/api/auth/complete", set.ApplicationHostname),
	}
}

func (a *AuthService) initStore(set *settings.Settings) {
	key := []byte(set.SessionSecret)
	store := sessions.NewCookieStore(key)
	store.Options.HttpOnly = true
	store.Options.Secure = true
	a.Store = store
}

func (a *AuthService) getState(req *http.Request) string {
	session := a.getSession(req)
	var state string
	stateval := session.Values[stateValue]
	if stateval != nil {
		state = stateval.(string)
	} else {
		nonceBytes := make([]byte, 64)
		io.ReadFull(rand.Reader, nonceBytes)
		state = base64.URLEncoding.EncodeToString(nonceBytes)
	}
	return state
}

func (a *AuthService) getSession(req *http.Request) *sessions.Session {
	session, err := a.Store.Get(req, sessionName)
	if err != nil {
		session, _ = a.Store.New(req, sessionName)
	}
	return session
}

// BeginAuthorization begins the OAuth2 flow.
func (a *AuthService) BeginAuthorization(res http.ResponseWriter, req *http.Request) {
	session := a.getSession(req)
	state := a.getState(req)
	session.Values[stateValue] = state
	session.Save(req, res)
	url := a.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// CompleteAuthorization attempts to complete the OAuth2 flow.
func (a *AuthService) CompleteAuthorization(res http.ResponseWriter, req *http.Request) {
	expectedState := a.getState(req)
	state := req.FormValue("state")
	if state != expectedState {
		http.Error(res, ErrBadState.Error(), http.StatusForbidden)
		return
	}
	code := req.FormValue("code")
	tok, err := a.OAuthConfig.Exchange(req.Context(), code)
	if err != nil {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}
	client := a.OAuthConfig.Client(req.Context(), tok)
	resp, err := client.Get("https://discordapp.com/api/users/@me")
	if err != nil {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}
	var du discordgo.User
	err = json.Unmarshal(body, &du)
	if err != nil {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseInt(du.ID, 10, 64)
	user := &domains.User{
		ID:   id,
		Name: fmt.Sprintf("%s#%s", du.Username, du.Discriminator),
	}
	session := a.getSession(req)
	userStr, err := json.Marshal(user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}
	session.Values[userValue] = userStr
	session.Save(req, res)
}

// GetAuthorization attempts to get the user's authorization object
func (a *AuthService) GetAuthorization(req *http.Request) (*domains.User, error) {
	session := a.getSession(req)
	var user domains.User
	if session != nil {
		userBytes := session.Values[userValue]
		if userBytes != nil {
			err := json.Unmarshal(userBytes.([]byte), &user)
			if err != nil {
				return nil, errors.Wrap(err, "could not unmarshal authorization")
			}
		}
	}
	return &user, nil
}

// IsAuthorized returns a simple boolean whether the user is authorized or not
func (a *AuthService) IsAuthorized(req *http.Request) bool {
	user, err := a.GetAuthorization(req)
	return err == nil && user.Name != ""
}

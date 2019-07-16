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
	"fmt"
	"github.com/kkragenbrink/slate/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type AuthSuite struct {
	suite.Suite
}

func TestAuthSuite(t *testing.T) {
	s := new(AuthSuite)
	suite.Run(t, s)
}

func newMockAuthSettings() *settings.Settings {
	set := new(settings.Settings)
	oauth := new(settings.OAuth)
	oauth.ClientID = "test"
	oauth.ClientSecret = "test"
	set.ApplicationHostname = "test"
	set.OAuth = oauth
	set.SessionSecret = "test"
	return set
}

func (suite *AuthSuite) TestNewAuthService() {
	a := NewAuthService(newMockAuthSettings())
	assert.NotNil(suite.T(), a)
}

func (suite *AuthSuite) TestBeginAuthorization() {
	a := NewAuthService(newMockAuthSettings())
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(suite.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(a.BeginAuthorization)
	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), http.StatusTemporaryRedirect, rr.Code)
	headers := rr.Header()
	cb := url.QueryEscape("https://test/api/auth/complete")
	state := url.QueryEscape(a.getState(req))
	expectedLocation := fmt.Sprintf("https://discordapp.com/api/oauth2/authorize?access_type=online&client_id=test&redirect_uri=%s&response_type=code&scope=identify&state=%s", cb, state)
	assert.Equal(suite.T(), expectedLocation, headers["Location"][0])
}

func (suite *AuthSuite) TestGetAuthorization() {
	a := NewAuthService(newMockAuthSettings())
	req, err := http.NewRequest("GET", "/", nil)
	sess := a.getSession(req)
	id := int64(1234)
	name := "test"
	bytes := fmt.Sprintf(`{"id":%d,"name":"%s"}`, id, name)
	sess.Values[userValue] = []byte(bytes)
	assert.Nil(suite.T(), err)
	user, err := a.GetAuthorization(req)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), id, user.ID)
	assert.Equal(suite.T(), name, user.Name)
}

func (suite *AuthSuite) TestIsAuthorized() {
	a := NewAuthService(newMockAuthSettings())
	req, err := http.NewRequest("GET", "/", nil)
	sess := a.getSession(req)
	id := int64(1234)
	name := "test"
	bytes := fmt.Sprintf(`{"id":%d,"name":"%s"}`, id, name)
	sess.Values[userValue] = []byte(bytes)
	assert.Nil(suite.T(), err)
	auth := a.IsAuthorized(req)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), auth)
}

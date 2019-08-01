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

package interfaces

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/go-chi/chi"

	"github.com/kkragenbrink/slate/domain"
)

type WebServiceTestSuite struct {
	suite.Suite

	now    time.Time
	router chi.Router
	rec    *httptest.ResponseRecorder
}

func (s *WebServiceTestSuite) SetupTest() {
	s.now = time.Now()
	cfg := &domain.SlateConfig{
		HerokuConfig: domain.HerokuConfig{
			ReleaseVersion:   "test",
			ReleaseCreatedAt: s.now.String(),
		},
	}
	s.router = SetupRoutes(cfg)
}

func (s *WebServiceTestSuite) TestVersion() {
	req, _ := http.NewRequest("GET", "/version", nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	body, err := ioutil.ReadAll(res.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var vr domain.VersionResponse
	err = json.Unmarshal(body, &vr)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), vr.Error)
	assert.Equal(s.T(), "test", vr.Response.Version)
	assert.Equal(s.T(), s.now.String(), vr.Response.ReleaseDate)
}

func TestWebService(t *testing.T) {
	suite.Run(t, new(WebServiceTestSuite))
}

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

package infrastructures

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kkragenbrink/slate/infrastructures/mocks"
	"github.com/kkragenbrink/slate/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type WebServiceSuite struct {
	suite.Suite
}

func TestWebServiceSuite(t *testing.T) {
	suite.Run(t, new(WebServiceSuite))
}

func (suite *WebServiceSuite) TestNewWebService() {
	set := new(settings.Settings)
	set.Port = 1234
	ws := NewWebService(set)
	assert.NotNil(suite.T(), ws)
}

func (suite *WebServiceSuite) TestHandlePost() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	expected := &TestBody{Content: "test"}
	jsonExpected, _ := json.Marshal(expected)
	router := mocks.NewMockRouter(ctrl)
	router.EXPECT().HandleFunc(gomock.Any(), gomock.Any()).Return(new(mux.Route))
	router.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
	ws := new(WebService)
	ws.router = router
	ws.HandlePost("/test", func(ctx context.Context, body json.RawMessage) (interface{}, error) {
		var got TestBody
		json.Unmarshal(body, &got)
		assert.Equal(suite.T(), expected, got)
		return got, nil
	})
	request, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonExpected))
	response := httptest.NewRecorder()
	ws.router.ServeHTTP(response, request)
	assert.Equal(suite.T(), 200, response.Code)
}

type TestBody struct {
	Content string `json:"content"`
}

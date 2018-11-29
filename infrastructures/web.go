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
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kkragenbrink/slate/interfaces"
	"github.com/kkragenbrink/slate/settings"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// A Router contains instructions for creating webservice handlers
type Router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// The WebService contains http-related services
type WebService struct {
	router   Router
	settings *settings.Settings
}

// NewWebService creates a new WebService instance and assigns a router to it.
func NewWebService(set *settings.Settings) *WebService {
	w := new(WebService)
	router := mux.NewRouter()
	w.router = router
	w.settings = set
	return w
}

// HandlePost handles a new JSON POST request
func (ws *WebService) HandlePost(route string, handler interfaces.WebHandler) {
	ws.router.HandleFunc(route, buildRequestHandler(handler)).Methods("POST")
}

// Start starts the web server and listener.
func (ws *WebService) Start() error {
	port := fmt.Sprintf(":%d", ws.settings.Port)
	err := http.ListenAndServe(port, ws.router)
	if err != nil {
		errors.Wrap(err, "could not establish web server")
	}
	return nil
}

func buildRequestHandler(handler interfaces.WebHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			// handle
			return
		}
		var jb json.RawMessage
		json.Unmarshal(body, &jb)

		re, err := handler(r.Context(), jb)
		if err != nil {
			w.WriteHeader(re.(int))
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		json.NewEncoder(w).Encode(re)
	}
}

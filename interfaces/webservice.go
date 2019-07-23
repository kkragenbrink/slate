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
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kkragenbrink/slate/domain"
	"net/http"
	"time"
)

// The WebServiceHandler contains handlers for http requests
type WebServiceHandler struct {
	Config *domain.SlateConfig
	Router *chi.Router
}

// SetupRoutes creates a chi router for use in handling HTTP requests
func SetupRoutes(cfg *domain.SlateConfig) chi.Router {
	ws := &WebServiceHandler{
		Config: cfg,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/version", ws.Version)

	return r
}

// Version handles requests to /version in the api
func (ws WebServiceHandler) Version(res http.ResponseWriter, req *http.Request) {
	response := domain.VersionResponse{}
	defer func() {
		if response.Error != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}
		output, _ := json.Marshal(response)
		res.Write(output)
	}()

	response.Response.Version = ws.Config.HerokuConfig.ReleaseVersion
	response.Response.ReleaseDate = ws.Config.HerokuConfig.ReleaseCreatedAt
}

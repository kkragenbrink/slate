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
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/kkragenbrink/slate/interfaces"
	"github.com/kkragenbrink/slate/settings"
	"net/http"
)

// The WebService manages the lifecycle of an http service
type WebService struct {
	handler  *interfaces.WebServiceHandler
	router   chi.Router
	server   *http.Server
	settings *settings.Settings
}

// NewWebService creates a new instance of the WebService manager
func NewWebService(set *settings.Settings, auth *AuthService, bot *Bot, db *DatabaseService, rand *RandomService) *WebService {
	ws := new(WebService)
	ws.router = initRouter()
	ws.server = initServer(set, ws.router)
	ws.handler = initHandler(ws.router, auth, bot, db, rand)
	return ws
}

func initHandler(router chi.Router, auth *AuthService, bot *Bot, db *DatabaseService, rand *RandomService) *interfaces.WebServiceHandler {
	handler := interfaces.NewWebServiceHandler(auth, bot, db, rand)
	router.Get("/auth", handler.Auth)
	router.Get("/auth/begin", handler.AuthBegin)
	router.Get("/auth/complete", handler.AuthComplete)
	router.Post("/channels", handler.Channels)
	router.Get("/characters", handler.Characters)
	router.Post("/roll", handler.Roll)
	router.Get("/sheets/{ID}", handler.Sheet)
	router.Post("/sheets/{ID}", handler.Sheet)
	return handler
}

func initRouter() chi.Router {
	router := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	router.Use(c.Handler)
	router.Use(middleware.Logger)
	return router
}

func initServer(set *settings.Settings, router chi.Router) *http.Server {
	server := new(http.Server)
	server.Addr = fmt.Sprintf(":%d", set.Port)
	server.Handler = router
	return server
}

// Start starts the http server listening
func (ws *WebService) Start() error {
	go ws.server.ListenAndServe()
	return nil
}

// Stop stops the http server from listening and shuts it down gracefully
func (ws *WebService) Stop() error {
	return ws.server.Shutdown(context.Background())
}

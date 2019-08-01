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

package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kkragenbrink/slate/interfaces"

	"github.com/kkragenbrink/slate/domain"

	"github.com/go-chi/chi"
)

// The HTTPServer controls interactions between our WebServices router and our HTTP server
type HTTPServer struct {
	cfg    *domain.SlateConfig
	router chi.Router
	run    *sync.Mutex
	server *http.Server
}

// NewHTTPServer initalizes the router and HTTP handlers, then returns an HTTPServer object
func NewHTTPServer(cfg *domain.SlateConfig) *HTTPServer {
	addr := fmt.Sprintf(":%d", cfg.Port)
	router := interfaces.SetupRoutes(cfg)
	http.Handle("/", router)

	return &HTTPServer{
		cfg:    cfg,
		router: router,
		run:    &sync.Mutex{},
		server: &http.Server{Addr: addr},
	}
}

// Start locks the mutex and then runs ListenAndServe in a goroutine
func (s *HTTPServer) Start() {
	s.run.Lock()
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			logrus.WithError(err).Fatal("could not start HTTP listener")
		}
	}()
}

// Stop attempts to shutdown the HTTP listener and then unlocks the mutex
func (s *HTTPServer) Stop() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		logrus.WithError(err).Error("could not gracefully shutdown HTTP listener")
	}
	s.run.Unlock()
}

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
	"github.com/kkragenbrink/slate/domain"
)

// Slate manages the state of the various application services
type Slate struct {
	config *domain.SlateConfig
	http   *HTTPServer
}

// NewSlate initializes the services and returns a new Slate instance
func NewSlate(cfg *domain.SlateConfig) *Slate {
	http := NewHTTPServer(cfg)

	slate := &Slate{
		config: cfg,
		http:   http,
	}
	return slate
}

// Start instructs the services to start
func (s *Slate) Start() {
	s.http.Start()
}

// Stop instructs the services to stop
func (s *Slate) Stop() {
	s.http.Stop()
}

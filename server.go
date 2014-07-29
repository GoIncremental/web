// Copyright 2014 GoIncremental Limited. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"net/http"
)

type Middleware interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type MiddlewareFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

type Server interface {
	Run(port string)
	Use(handler Middleware)
	UseHandler(http.Handler)
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
}

type server struct {
	negroni *negroni.Negroni
}

func (s *server) Run(port string) {
	s.negroni.Run(port)
}

func (s *server) Use(handler Middleware) {
	s.negroni.Use(handler)
}

func (s *server) UseHandler(h http.Handler) {
	s.negroni.UseHandler(h)
}

func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.negroni.ServeHTTP(rw, r)
}

func NewServer() Server {
	n := negroni.Classic()
	return &server{negroni: n}
}

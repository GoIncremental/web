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
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	http.Handler
	HandleFunc(s string, f func(http.ResponseWriter, *http.Request)) Route
	Handle(path string, handler http.Handler) Route
}

type Route interface {
	Methods(s ...string) Route
}

type route struct {
	Route
	route *mux.Route
}

func (r *route) Methods(s ...string) Route {
	muxRoute := r.route.Methods(s...)
	r.route = muxRoute
	return r
}

type router struct {
	Router
	router *mux.Router
}

func (r *router) HandleFunc(s string, f func(http.ResponseWriter, *http.Request)) Route {
	muxRoute := r.router.HandleFunc(s, f)
	return &route{route: muxRoute}
}

func (r *router) Handle(path string, handler http.Handler) Route {
	muxRoute := r.router.Handle(path, handler)
	return &route{route: muxRoute}
}

func (r *router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(rw, req)
}

func NewRouter() Router {
	m := mux.NewRouter()
	g := &router{router: m}
	return g
}

func Params(req *http.Request) map[string]string {
	return mux.Vars(req)
}

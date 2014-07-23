package web

import (
	"github.com/gorilla/mux"
	"net/http"
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

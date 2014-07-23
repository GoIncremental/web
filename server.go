package web

import (
	"github.com/codegangsta/negroni"
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

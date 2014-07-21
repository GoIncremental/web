package web

import (
	"github.com/codegangsta/negroni"
	"net/http"
)

type Middleware interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type Server interface {
	Run(port string)
	Use(handler Middleware)
	UseHandler(http.Handler)
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

func NewServer() Server {
	n := negroni.Classic()
	return &server{negroni: n}
}

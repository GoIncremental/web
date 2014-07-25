package web

import (
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/goincremental/dal"
	"net/http"
)

type Session struct {
	ID      string
	Values  map[interface{}]interface{}
	Options *Options
	IsNew   bool
	name    string
	written bool
}

func (s *Session) Get(key interface{}) interface{} {
	return s.Values[key]
}

func (s *Session) Set(key interface{}, val interface{}) {
	s.Values[key] = val
	s.written = true
}

func (s *Session) Delete(key interface{}) {
	delete(s.Values, key)
	s.written = true
}

type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// NewSessionStore returns a new SessionStore (currently uses default dal implementation)
// Set ensureTTL to true let the database auto-remove expired object by maxAge.
func NewSessionStore(c dal.Collection, maxAge int, ensureTTL bool, keyPairs ...[]byte) SessionStore {
	return newDalStore(c, maxAge, ensureTTL, keyPairs...)
}

type SessionStore interface {
	GetSession(r *http.Request, name string) (s *Session, err error)
	SaveSession(r *http.Request, w http.ResponseWriter, session *Session) error
}

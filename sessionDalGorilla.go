package web

import (
	"errors"
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/goincremental/dal"
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/gorilla/securecookie"
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/gorilla/sessions"
	"net/http"
	"time"
)

var (
	ErrInvalidId       = errors.New("session: invalid session id")
	ErrInvalidModified = errors.New("mongostore: invalid modified value")
)

type tokenGetSeter interface {
	getToken(req *http.Request, name string) (string, error)
	setToken(rw http.ResponseWriter, name, value string, options *sessions.Options)
}

type cookieToken struct{}

func (c *cookieToken) getToken(req *http.Request, name string) (string, error) {
	cook, err := req.Cookie(name)
	if err != nil {
		return "", err
	}

	return cook.Value, nil
}

func (c *cookieToken) setToken(rw http.ResponseWriter, name string, value string,
	options *sessions.Options) {
	http.SetCookie(rw, sessions.NewCookie(name, value, options))
}

func gorilla2session(gs *sessions.Session) *Session {
	return &Session{
		ID:     gs.ID,
		Values: gs.Values,
		Options: &Options{
			Path:     gs.Options.Path,
			Domain:   gs.Options.Domain,
			MaxAge:   gs.Options.MaxAge,
			Secure:   gs.Options.Secure,
			HttpOnly: gs.Options.HttpOnly,
		},
		IsNew: gs.IsNew,
		name:  gs.Name(),
	}
}

// Session object store via dal
type dalSession struct {
	Id       dal.ObjectId `bson:"_id,omitempty"`
	Data     string
	Modified time.Time
}

func newDalStore(c dal.Collection, maxAge int, ensureTTL bool, keyPairs ...[]byte) *dalStore {
	if ensureTTL {
		c.EnsureIndex(dal.Index{
			Key:         []string{"modified"},
			Background:  true,
			Sparse:      true,
			ExpireAfter: time.Duration(maxAge) * time.Second,
		})
	}
	return &dalStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			MaxAge: maxAge,
		},
		Token: &cookieToken{},
		coll:  c,
	}
}

type dalStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options
	Token   tokenGetSeter
	coll    dal.Collection
}

func (m *dalStore) fromSession(s *Session) *sessions.Session {
	gsess := sessions.NewSession(m, s.name)
	gsess.ID = s.ID
	gsess.Values = s.Values
	gsess.Options = m.Options
	gsess.IsNew = s.IsNew
	return gsess
}

//Implementation of web.SessionStore Interface
func (m *dalStore) GetSession(r *http.Request, name string) (s *Session, err error) {
	gsess, err := m.Get(r, name)
	if err == nil {
		s = gorilla2session(gsess)
	}
	return
}

func (m *dalStore) SaveSession(r *http.Request, w http.ResponseWriter, session *Session) error {
	return m.Save(r, w, m.fromSession(session))
}

//Implementation of gorilla/sessions.Store interface
// Get registers and returns a session for the given name and session store.
// It returns a new session if there are no sessions registered for the name.
func (m *dalStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(m, name)
}

// New returns a session for the given name without adding it to the registry.
func (m *dalStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(m, name)
	session.Options = &sessions.Options{
		Path:   m.Options.Path,
		MaxAge: m.Options.MaxAge,
	}
	session.IsNew = true
	var err error
	if cook, errToken := m.Token.getToken(r, name); errToken == nil {
		err = securecookie.DecodeMulti(name, cook, &session.ID, m.Codecs...)
		if err == nil {
			err = m.load(session)
			if err == nil {
				session.IsNew = false
			} else {
				err = nil
			}
		}
	}
	return session, err
}

func (m *dalStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge < 0 {
		if err := m.delete(session); err != nil {
			return err
		}
		m.Token.setToken(w, session.Name(), "", session.Options)
		return nil
	}

	if session.ID == "" {
		session.ID = dal.NewObjectId().Hex()
	}

	if err := m.upsert(session); err != nil {
		return err
	}

	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
		m.Codecs...)
	if err != nil {
		return err
	}

	m.Token.setToken(w, session.Name(), encoded, session.Options)
	return nil
}

func (m *dalStore) load(session *sessions.Session) error {
	if !dal.IsObjectIdHex(session.ID) {
		return ErrInvalidId
	}

	s := dalSession{}
	err := m.coll.FindId(dal.ObjectIdHex(session.ID)).One(&s)
	if err != nil {
		return err
	}

	if err := securecookie.DecodeMulti(session.Name(), s.Data, &session.Values,
		m.Codecs...); err != nil {
		return err
	}

	return nil
}

func (m *dalStore) upsert(session *sessions.Session) error {
	if !dal.IsObjectIdHex(session.ID) {
		return ErrInvalidId
	}

	var modified time.Time
	if val, ok := session.Values["modified"]; ok {
		modified, ok = val.(time.Time)
		if !ok {
			return ErrInvalidModified
		}
	} else {
		modified = time.Now()
	}

	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		m.Codecs...)
	if err != nil {
		return err
	}

	s := dalSession{
		Id:       dal.ObjectIdHex(session.ID),
		Data:     encoded,
		Modified: modified,
	}

	_, err = m.coll.UpsertId(s.Id, &s)
	if err != nil {
		return err
	}

	return nil
}

func (m *dalStore) delete(session *sessions.Session) error {
	if !dal.IsObjectIdHex(session.ID) {
		return ErrInvalidId
	}

	return m.coll.RemoveId(dal.ObjectIdHex(session.ID))
}

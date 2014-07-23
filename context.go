package web

import (
	"github.com/goincremental/dal"
	"github.com/gorilla/context"
	"net/http"
)

func getContext(r *http.Request, key interface{}) interface{} {
	return context.Get(r, key)
}

func setContext(r *http.Request, key, val interface{}) {
	context.Set(r, key, val)
}

type contextKey int

const dbKey contextKey = 0
const rendererKey contextKey = 4
const sessionKey contextKey = 5
const sessionStoreKey contextKey = 6

func GetDb(r *http.Request) dal.Database {
	if rv := getContext(r, dbKey); rv != nil {
		return rv.(dal.Database)
	}
	return nil
}

func SetDb(r *http.Request, val dal.Database) {
	setContext(r, dbKey, val)
}

func GetRenderer(r *http.Request) Renderer {
	if rv := getContext(r, rendererKey); rv != nil {
		return rv.(Renderer)
	}
	return nil
}

func SetRenderer(r *http.Request, val Renderer) {
	setContext(r, rendererKey, val)
}

func GetSession(r *http.Request) *Session {
	if rv := getContext(r, sessionKey); rv != nil {
		return rv.(*Session)
	}
	return nil
}

func SetSession(r *http.Request, val *Session) {
	setContext(r, sessionKey, val)
}

func GetSessionStore(r *http.Request) SessionStore {
	if rv := getContext(r, sessionStoreKey); rv != nil {
		return rv.(SessionStore)
	}
	return nil
}

func SetSessionStore(r *http.Request, val SessionStore) {
	setContext(r, sessionStoreKey, val)
}

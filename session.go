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
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/goincremental/dal"
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/goincremental/negroni-sessions"
	"net/http"
)

type Session interface {
	sessions.Session
}

// Store is an interface for custom session stores.
type Store interface {
	sessions.Store
}

func Sessions(name string, store Store) MiddlewareFunc {
	return MiddlewareFunc(sessions.Sessions(name, store))
}

func GetSession(req *http.Request) Session {
	return sessions.GetSession(req)
}

// NewSessionStore returns a new SessionStore (currently uses default dal implementation)
// Set ensureTTL to true let the database auto-remove expired object by maxAge.
func NewSessionStore(c dal.Collection, maxAge int, ensureTTL bool, keyPairs ...[]byte) Store {
	return sessions.NewDalStore(c, maxAge, ensureTTL, keyPairs...)
}

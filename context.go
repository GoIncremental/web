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

	"github.com/goincremental/dal"
	"github.com/gorilla/context"
)

func GetContext(r *http.Request, key interface{}) interface{} {
	return context.Get(r, key)
}

func SetContext(r *http.Request, key, val interface{}) {
	context.Set(r, key, val)
}

type contextKey int

const dbKey contextKey = 0
const rendererKey contextKey = 4
const sessionKey contextKey = 5
const sessionStoreKey contextKey = 6

func GetDb(r *http.Request) *dal.Database {
	if rv := GetContext(r, dbKey); rv != nil {
		db := rv.(dal.Database)
		return &db
	}
	return nil
}

func SetDb(r *http.Request, val dal.Database) {
	SetContext(r, dbKey, val)
}

func GetRenderer(r *http.Request) Renderer {
	if rv := GetContext(r, rendererKey); rv != nil {
		return rv.(Renderer)
	}
	return nil
}

func SetRenderer(r *http.Request, val Renderer) {
	SetContext(r, rendererKey, val)
}

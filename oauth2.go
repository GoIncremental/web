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

// Package oauth2 contains Negroni middleware to provide
// user login via an OAuth 2.0 backend.
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/goincremental/web/Godeps/_workspace/src/github.com/golang/oauth2"
)

const (
	keyToken    = "oauth2_token"
	keyNextPage = "next"
)

var (
	// Path to handle OAuth 2.0 logins.
	PathLogin = "/login"
	// Path to handle OAuth 2.0 logouts.
	PathLogout = "/logout"
	// Path to handle callback from OAuth 2.0 backend
	// to exchange credentials.
	PathCallback = "/oauth2callback"
	// Path to handle error cases.
	PathError = "/oauth2error"
)

// Represents a container that contains
// user's OAuth 2.0 access and refresh tokens.
type Tokens interface {
	Access() string
	Refresh() string
	IsExpired() bool
	ExpiryTime() time.Time
	ExtraData() map[string]string
}

type token struct {
	oauth2.Token
}

func (t *token) ExtraData() map[string]string {
	return t.Extra
}

// Returns the access token.
func (t *token) Access() string {
	return t.AccessToken
}

// Returns the refresh token.
func (t *token) Refresh() string {
	return t.RefreshToken
}

// Returns whether the access token is
// expired or not.
func (t *token) IsExpired() bool {
	if t == nil {
		return true
	}
	return t.Expired()
}

// Returns the expiry time of the user's
// access token.
func (t *token) ExpiryTime() time.Time {
	return t.Expiry
}

// // Formats tokens into string.
// func (t *token) String() string {
// 	return fmt.Sprintf("tokens: %v", t)
// }

// Returns a new Google OAuth 2.0 backend endpoint.
func Google(opts *oauth2.Options) http.Handler {
	authUrl := "https://accounts.google.com/o/oauth2/auth"
	tokenUrl := "https://accounts.google.com/o/oauth2/token"
	return NewOAuth2Provider(opts, authUrl, tokenUrl)
}

// Returns a new Github OAuth 2.0 backend endpoint.
func Github(opts *oauth2.Options) http.Handler {
	authUrl := "https://github.com/login/oauth/authorize"
	tokenUrl := "https://github.com/login/oauth/access_token"
	return NewOAuth2Provider(opts, authUrl, tokenUrl)
}

func Facebook(opts *oauth2.Options) http.Handler {
	authUrl := "https://www.facebook.com/dialog/oauth"
	tokenUrl := "https://graph.facebook.com/oauth/access_token"
	return NewOAuth2Provider(opts, authUrl, tokenUrl)
}

func LinkedIn(opts *oauth2.Options) http.Handler {
	authUrl := "https://www.linkedin.com/uas/oauth2/authorization"
	tokenUrl := "https://www.linkedin.com/uas/oauth2/accessToken"
	return NewOAuth2Provider(opts, authUrl, tokenUrl)
}

// Returns a generic OAuth 2.0 backend endpoint.
func NewOAuth2Provider(opts *oauth2.Options, authUrl, tokenUrl string) http.HandlerFunc {

	config, err := oauth2.NewConfig(opts, authUrl, tokenUrl)
	if err != nil {
		panic(fmt.Sprintf("oauth2: %s", err))
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		s := GetSession(r)
		store := GetSessionStore(r)

		if r.Method == "GET" {
			switch r.URL.Path {
			case PathLogin:
				login(config, s, rw, r)
			case PathLogout:
				logout(s, rw, r)
			case PathCallback:
				handleOAuth2Callback(config, s, rw, r)
			}
		}

		tk := unmarshallToken(*s)
		if tk != nil {
			// check if the access token is expired
			if tk.IsExpired() && tk.Refresh() == "" {
				s.Delete(keyToken)
				store.SaveSession(r, rw, s)
				tk = nil
			}
		}
	}
}

// Handler that redirects user to the login page
// if user is not logged in.
func LoginRequired() Middleware {
	return MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s := GetSession(r)
		token := unmarshallToken(*s)
		if token == nil || token.IsExpired() {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(rw, r, PathLogin+"?next="+next, http.StatusFound)
		} else {
			next(rw, r)
		}
	})
}

func login(c *oauth2.Config, s *Session, w http.ResponseWriter, r *http.Request) {
	next := extractPath(r.URL.Query().Get(keyNextPage))
	if s.Get(keyToken) == nil {
		// User is not logged in.
		if next == "" {
			next = "/"
		}
		http.Redirect(w, r, c.AuthCodeURL(next), http.StatusFound)
		return
	}
	// No need to login, redirect to the next page.
	http.Redirect(w, r, next, http.StatusFound)
}

func logout(s *Session, w http.ResponseWriter, r *http.Request) {
	next := extractPath(r.URL.Query().Get(keyNextPage))
	s.Delete(keyToken)
	http.Redirect(w, r, next, http.StatusFound)
}

func handleOAuth2Callback(c *oauth2.Config, s *Session, w http.ResponseWriter, r *http.Request) {
	next := extractPath(r.URL.Query().Get("state"))
	code := r.URL.Query().Get("code")
	t, err := c.NewTransportWithCode(code)
	if err != nil {
		// Pass the error message, or allow dev to provide its own
		// error handler.
		http.Redirect(w, r, PathError, http.StatusFound)
		return
	}
	// Store the credentials in the session.
	val, _ := json.Marshal(t.Token())
	s.Set(keyToken, val)
	store := GetSessionStore(r)
	store.SaveSession(r, w, s)
	http.Redirect(w, r, next, http.StatusFound)
}

func unmarshallToken(s Session) (t *token) {

	if s.Values[keyToken] == nil {
		return
	}
	data := s.Values[keyToken].([]byte)
	var tk oauth2.Token
	json.Unmarshal(data, &tk)
	return &token{tk}
}

func extractPath(next string) string {
	n, err := url.Parse(next)
	if err != nil {
		return "/"
	}
	return n.Path
}

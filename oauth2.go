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
	"time"

	oauth2 "github.com/goincremental/negroni-oauth2"
)

// Returns a new Google OAuth 2.0 backend endpoint.
func GoogleOLD(opts *OAuth2Options) Middleware {
	authUrl := "https://accounts.google.com/o/oauth2/auth"
	tokenUrl := "https://accounts.google.com/o/oauth2/token"
	return NewOAuth2Provider(opts, authUrl, tokenUrl)
}

// Returns a generic OAuth 2.0 backend endpoint.
func NewOAuth2Provider(opts *OAuth2Options, authUrl, tokenUrl string) Middleware {
	options := &oauth2.Config{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURL:  opts.RedirectURL,
		Scopes:       opts.Scopes,
	}
	return oauth2.NewOAuth2Provider(options, authUrl, tokenUrl)
}

// Handler that redirects user to the login page
// if user is not logged in.
func LoginRequired() Middleware {
	return MiddlewareFunc(oauth2.LoginRequired())
}

type OAuthToken interface {
	Access() string
	Refresh() string
	Expired() bool
	ExpiryTime() time.Time
	ExtraData(string) string
}

func GetOAuth2Token(r *http.Request) OAuthToken {
	return oauth2.GetToken(r)
}

type OAuth2Options struct {
	// ClientID is the OAuth client identifier used when communicating with
	// the configured OAuth provider.
	ClientID string `json:"client_id"`

	// ClientSecret is the OAuth client secret used when communicating with
	// the configured OAuth provider.
	ClientSecret string `json:"client_secret"`

	// RedirectURL is the URL to which the user will be returned after
	// granting (or denying) access.
	RedirectURL string `json:"redirect_url"`

	// Optional, identifies the level of access being requested.
	Scopes []string `json:"scopes"`
}

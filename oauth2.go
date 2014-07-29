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
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/goincremental/negroni-oauth2"
)

// Returns a new Google OAuth 2.0 backend endpoint.
func Google(opts *oauth2.Options) Middleware {
	return oauth2.Google(opts)
}

// Returns a new Github OAuth 2.0 backend endpoint.
func Github(opts *oauth2.Options) Middleware {
	return oauth2.Github(opts)
}

func Facebook(opts *oauth2.Options) Middleware {
	return oauth2.Facebook(opts)
}

func LinkedIn(opts *oauth2.Options) Middleware {
	return oauth2.LinkedIn(opts)
}

// Returns a generic OAuth 2.0 backend endpoint.
func NewOAuth2Provider(opts *oauth2.Options, authUrl, tokenUrl string) Middleware {
	return oauth2.NewOAuth2Provider(opts, authUrl, tokenUrl)
}

// Handler that redirects user to the login page
// if user is not logged in.
func LoginRequired() Middleware {
	return MiddlewareFunc(oauth2.LoginRequired())
}

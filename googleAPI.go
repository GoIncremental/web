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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	plus "code.google.com/p/google-api-go-client/plus/v1"
	"github.com/golang/oauth2"
)

// GoogleAPI supports quering the google api to get login information about
// the current user, and call other google APIs on behalf of that user
type GoogleAPI interface {
	LoginWithCode(code string) (GoogleUser, error)
}

type googleAPI struct {
	config *oauth2.Config
}

func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

// claimSet represents an IdToken response.
type claimSet struct {
	Sub string
}

// decodeIdToken takes an ID Token and decodes it to fetch the Google+ ID within
func decodeToken(idToken string) (gplusID string, err error) {
	// An ID token is a cryptographically-signed JSON object encoded in base 64.
	// Normally, it is critical that you validate an ID token before you use it,
	// but since you are communicating directly with Google over an
	// intermediary-free HTTPS channel and using your Client Secret to
	// authenticate yourself to Google, you can be confident that the token you
	// receive really comes from Google and is valid. If your server passes the ID
	// token to other components of your app, it is extremely important that the
	// other components validate the token before using it.
	var set claimSet
	if idToken != "" {
		// Check that the padding is correct for a base64decode
		parts := strings.Split(idToken, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("Malformed ID token")
		}
		// Decode the ID token
		b, err := base64Decode(parts[1])
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
		err = json.Unmarshal(b, &set)
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
	}
	return set.Sub, nil
}

func (g *googleAPI) LoginWithCode(code string) (GoogleUser, error) {
	var valid = false

	token, err := g.config.Exchange(nil, code)
	if err != nil {
		return nil, err
	}

	id, err := decodeToken(token.Extra("id_token"))
	if err != nil {
		return nil, err
	}

	client := g.config.Client(oauth2.NoContext, token)
	service, err := plus.New(client)
	getme := service.People.Get("me")
	me, err := getme.Do()

	email := "N/A"
	for _, e := range me.Emails {
		if e.Type == "account" {
			email = e.Value
		}
	}
	valid = true
	return &googleUser{id: id, valid: valid, email: email}, nil
}

// GoogleUser provides information about the currently logged in user acording
// to the google api
type GoogleUser interface {
	ID() string
	Email() string
	Valid() bool
}

type googleUser struct {
	id    string
	email string
	valid bool
}

func (g *googleUser) ID() string {
	return g.id
}

func (g *googleUser) Email() string {
	return g.email
}

func (g *googleUser) Valid() bool {
	return g.valid
}

// NewGoogleAPI creates a GoogleAPI object.  This should be created
// once per web server, and then stored on the request using the SetGoogle
// middleware
func NewGoogleAPI(id string, secret string, scopes []string) GoogleAPI {
	authURL := "https://accounts.google.com/o/oauth2/auth"
	tokenURL := "https://accounts.google.com/o/oauth2/token"
	config := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  "postmessage",
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		}}

	return &googleAPI{config: config}
}

const googleAPIKey contextKey = 7

// SetGoogleAPI is middleware that ensures the provided GoogleAPI object
// is made available on the request object
func SetGoogleAPI(googleAPI GoogleAPI) Middleware {
	return MiddlewareFunc(func(rw http.ResponseWriter,
		r *http.Request, next http.HandlerFunc) {
		SetContext(r, googleAPIKey, googleAPI)
		next(rw, r)
	})
}

// GetGoogleAPI retrieves the GoogleAPI object from the request context
func GetGoogleAPI(r *http.Request) GoogleAPI {
	if rv := GetContext(r, googleAPIKey); rv != nil {
		return rv.(GoogleAPI)
	}
	return nil
}

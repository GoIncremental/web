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
	"github.com/goincremental/web/Godeps/_workspace/src/github.com/joho/godotenv"
	"log"
)

type Environment interface {
	Load(filenames ...string) error
}

type environment struct{}

func (e *environment) Load(s ...string) (err error) {
	err = godotenv.Load(s...)
	return
}

func newEnvironment() Environment {
	return &environment{}
}

func LoadEnv() {
	env := newEnvironment()
	err := env.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

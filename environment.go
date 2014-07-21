package web

import (
	"github.com/joho/godotenv"
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

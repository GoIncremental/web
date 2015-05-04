package web

import (
	"fmt"
	"github.com/deferpanic/deferclient/deferstats"
	"github.com/deferpanic/deferclient/errors"
	"log"
	"os"
)

// NewLogger configures defer panic error and starts capturing stats.
// It looks for DEFERPANIC_KEY,
// DEVERPANIC_ENVIRONMENT and DEFERPANIC_APPGROUP environment vars
func NewLogger() {
	dfs := deferstats.NewClient(os.Getenv("DEFERPANIC_KEY"))
	dfs.Setenvironment(os.Getenv("DEFERPANIC_ENVIRONMENT"))
	dfs.SetappGroup(os.Getenv("DEFERPANIC_APPGROUP"))
	go dfs.CaptureStats()
}

func logError(msg string) {
	errors.New(msg)
	log.Println(msg)
}

// LogError passes the error to deferpanic
// and prints the error message to stdout
func LogError(e error) {
	logError(e.Error())
}

// LogErrorf accepts a format string and arguments
// It creates a new error, logs with deferpanic and
// prints the error to stdout
func LogErrorf(format string, a ...interface{}) {
	LogError(fmt.Errorf(format, a...))
}

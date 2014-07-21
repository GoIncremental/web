package web

import (
	"github.com/gorilla/context"
	"net/http"
)

func GetContext(r *http.Request, key interface{}) interface{} {
	return context.Get(r, key)
}

func SetContext(r *http.Request, key, val interface{}) {
	context.Set(r, key, val)
}

package security

import (
	"net/http"

	"github.com/goincremental/web"
	"github.com/goincremental/web/security/models"
)

type securityContextKey int

const userKey securityContextKey = 0

func SetUser(r *http.Request, val *models.User) {
	web.SetContext(r, userKey, val)
}

func GetUser(r *http.Request) *models.User {
	if u := web.GetContext(r, userKey); u != nil {
		user := u.(*models.User)
		return user
	}
	return nil
}

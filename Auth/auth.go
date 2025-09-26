package Auth

import (
	AuthInterface "dependencies/Auth/Interface"
	"net/http"
)

type Authenticate struct {
	AuthMiddleware AuthInterface.IAuthMiddleware
}

func NewAuthenticate(AuthMiddleware AuthInterface.IAuthMiddleware) Authenticate {
	authenticate := Authenticate{
		AuthMiddleware: AuthMiddleware,
	}
	return authenticate
}

func (a Authenticate) authenticate(next http.Handler) {
	a.AuthMiddleware.Authenticate(next)
}

package Interface

import "net/http"

type IAuthMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}

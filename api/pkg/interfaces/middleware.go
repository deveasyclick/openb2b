package interfaces

import (
	"net/http"

	clerkHttp "github.com/clerk/clerk-sdk-go/v2/http"
)

type Middleware interface {
	Recover(logger Logger) func(http.Handler) http.Handler
	ValidateJWT(opts ...clerkHttp.AuthorizationOption) func(http.Handler) http.Handler
}

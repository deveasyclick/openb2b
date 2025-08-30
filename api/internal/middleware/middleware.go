package middleware

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkHttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type middleware struct{}

func New() interfaces.Middleware {
	return &middleware{}
}

func customClaimsConstructor(ctx context.Context) any {
	return &identity.CustomSessionClaims{}
}

func withCustomClaimsConstructor(params *clerkHttp.AuthorizationParams) error {
	params.VerifyParams.CustomClaimsConstructor = customClaimsConstructor
	return nil
}

// ValidateJWT middleware enforces Clerk authentication on protected routes.
//
// This middleware:
//  1. Wraps the Clerk `WithHeaderAuthorization` middleware to validate
//     incoming requests using the `Authorization` header.
//  2. Attaches custom claims extraction via `WithCustomClaimsConstructor`
//     so application-specific claims can be included in the session context.
//  3. Retrieves Clerk session claims from the request context.
//  4. If no valid session is found, responds with `401 Unauthorized` and
//     a JSON error message.
//  5. If authentication succeeds, passes control to the next handler.
//
// Parameters:
//   - opts ...clerkHttp.AuthorizationOption: optional Clerk authorization options
//     (e.g., audience, allowed origins, or custom claim configuration).
//
// Usage (Chi example):
//
//	r := chi.NewRouter()
//	r.Group(func(r chi.Router) {
//	    r.Use(middleware.ValidateJWT()) // Protect all routes in this group
//	    r.Get("/profile", profileHandler)
//	})
//
// Behavior:
//   - Requests with a valid Clerk session in the `Authorization` header will
//     continue to the next handler.
//   - Requests without a valid session will be rejected with 401 Unauthorized.
func (m *middleware) ValidateJWT(opts ...clerkHttp.AuthorizationOption) func(http.Handler) http.Handler {
	// Extract custom claims from the request
	opts = append(opts, withCustomClaimsConstructor)
	return func(next http.Handler) http.Handler {
		return clerkHttp.WithHeaderAuthorization(opts...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok || claims == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "unauthorized"}`))
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

func (m *middleware) Recover(logger interfaces.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered", "error", rec)
					response.WriteJSONError(w, &apperrors.APIError{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error",
					}, logger)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

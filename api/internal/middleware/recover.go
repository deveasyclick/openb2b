package middleware

import (
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

// Recover returns a middleware that recovers from panics.
// Only panics are logged with a stack trace.
func Recover(logger interfaces.Logger) func(http.Handler) http.Handler {
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

package setup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkHttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type fakeMiddleware struct {
	UserID  uint
	OrgID   uint
	ClerkID string
}

func NewFake(userId uint, orgId uint, clerkId string) interfaces.Middleware {
	return &fakeMiddleware{
		UserID:  1,
		OrgID:   1,
		ClerkID: clerkId,
	}
}

func (m *fakeMiddleware) Recover(logger interfaces.Logger) func(http.Handler) http.Handler {
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

// ValidateJWT in tests simply injects a fake Clerk session into context.
func (m *fakeMiddleware) ValidateJWT(opts ...clerkHttp.AuthorizationOption) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create fake claims
			fakeClaims := &clerk.SessionClaims{
				RegisteredClaims: clerk.RegisteredClaims{
					Subject: m.ClerkID,
				},
				Custom: &identity.CustomSessionClaims{
					UserID: fmt.Sprintf("%d", m.UserID),
					OrgID:  fmt.Sprintf("%d", m.OrgID),
				},
			}

			// Put them into context (mimicking Clerk)
			ctx := clerk.ContextWithSessionClaims(r.Context(), fakeClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *fakeMiddleware) VerifyWebhook() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Error reading request body"}`))
				return
			}
			defer r.Body.Close()

			// Create a new reader with the body for verification
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			event := types.WebhookEvent{}
			if err := json.Unmarshal(body, &event); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Invalid webhook payload"}`))
				return
			}
			// Store the parsed event in the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, types.WebhookEventKey, event)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

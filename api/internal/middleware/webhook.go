package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/svix"
)

// Verify is a middleware that validates incoming webhook requests.
//
// It performs the following steps:
//  1. Initializes the Svix webhook verifier using the configured Clerk secret.
//  2. Reads and buffers the request body so it can be used for both verification
//     and downstream handlers.
//  3. Extracts the required `svix-id`, `svix-timestamp`, and `svix-signature`
//     headers from the incoming request.
//  4. Verifies the request signature against the Clerk secret to ensure the
//     payload is authentic and untampered.
//  5. Parses the request body into a `types.WebhookEvent` struct.
//  6. Stores the parsed event inside the request context under `webhookEventKey`
//     so that downstream handlers (e.g., webhook handlers) can access it directly.
//
// On error (invalid signature, bad body, or malformed payload), this middleware
// writes a JSON error response and terminates the request.
//
// Usage (Chi example):
//
//	r := chi.NewRouter()
//	appCtx := deps.NewAppContext()
//	r.Use(webhook.VerifyWebhook(appCtx))
//	r.Post("/webhook/clerk", handler.ClerkWebhookHandler)
//
// Dependencies:
//   - appCtx.Config.ClerkWebhookSigningSecret: the secret used to validate signatures.
//   - appCtx.Logger: for structured error logging.
//
// Example flow:
//
//	Clerk → POST /webhook/clerk → VerifyWebhook middleware (validates + injects event)
//	→ ClerkWebhookHandler (consumes types.WebhookEvent from context).
func (m *middleware) VerifyWebhook() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify webhook signature
			wh, err := svix.GetWebhookVerifier(m.appCtx.Config.ClerkWebhookSigningSecret)
			if err != nil {
				response.WriteJSONError(w, &apperrors.APIError{
					Code:    http.StatusInternalServerError,
					Message: "Webhook verifier not initialized",
				}, m.appCtx.Logger)
				return
			}

			// Read and validate the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				response.WriteJSONError(w, &apperrors.APIError{
					Code:    http.StatusBadRequest,
					Message: "Error reading request body",
				}, m.appCtx.Logger)
				return
			}
			defer r.Body.Close()

			// Create a new reader with the body for verification
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Create headers for verification
			headers := http.Header{}
			headers.Set("svix-id", r.Header.Get("svix-id"))
			headers.Set("svix-timestamp", r.Header.Get("svix-timestamp"))
			headers.Set("svix-signature", r.Header.Get("svix-signature"))

			// Verify the webhook
			err = wh.Verify(body, headers)
			if err != nil {
				response.WriteJSONError(w, &apperrors.APIError{
					Code:    http.StatusUnauthorized,
					Message: "Invalid webhook signature",
				}, m.appCtx.Logger)
				return
			}

			// Parse the webhook event
			event := types.WebhookEvent{}
			if err := json.Unmarshal(body, &event); err != nil {
				response.WriteJSONError(w, &apperrors.APIError{
					Code:    http.StatusBadRequest,
					Message: "Invalid webhook payload",
				}, m.appCtx.Logger)
				return
			}
			// Store the parsed event in the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, types.WebhookEventKey, event)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

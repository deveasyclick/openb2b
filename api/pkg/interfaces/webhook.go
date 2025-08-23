package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/types"
)

type WebhookService interface {
	HandleEvent(ctx context.Context, event *types.WebhookEvent) *apperrors.APIError
}

type WebhookHandler interface {
	ClerkWehbookHandler(w http.ResponseWriter, r *http.Request)
}

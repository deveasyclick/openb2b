package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/types"
)

type WebhookService interface {
	HandleEvent(ctx context.Context, event *types.WebhookEvent) error
}

type WebhookHandler interface {
	HandleClerkEvents(w http.ResponseWriter, r *http.Request)
}

package webhook

import (
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type webhookHandler struct {
	webhookService interfaces.WebhookService
	appCtx         *deps.AppContext
}

func NewWebhookHandler(webhookService interfaces.WebhookService, appCtx *deps.AppContext) interfaces.WebhookHandler {
	return &webhookHandler{
		webhookService: webhookService,
		appCtx:         appCtx,
	}
}

func (h *webhookHandler) ClerkWehbookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event := ctx.Value(webhookEventKey).(*types.WebhookEvent)
	if event == nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: "No webhook event in context",
		}, h.appCtx.Logger)
		return
	}

	h.appCtx.Logger.Info("Handling webhook event", "event", event.Data)

	err := h.webhookService.HandleEvent(ctx, event)

	if err != nil {
		response.WriteJSONError(w, err, h.appCtx.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package webhook

import (
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type handler struct {
	webhookService interfaces.WebhookService
	appCtx         *deps.AppContext
}

func NewHandler(webhookService interfaces.WebhookService, appCtx *deps.AppContext) interfaces.WebhookHandler {
	return &handler{
		webhookService: webhookService,
		appCtx:         appCtx,
	}
}

// HandleClerkEvents godoc
// @Summary      Receive Clerk webhook events
// @Description  Handles incoming webhook events from Clerk. The request body must match the WebhookEvent structure.
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        event  body      types.WebhookEvent  true  "Webhook Event Payload"
// @Success      200    {string}  string               "OK"
// @Failure      400    {object}  apperrors.APIErrorResponse
// @Failure      401    {object}  apperrors.APIErrorResponse
// @Failure      500    {object}  apperrors.APIErrorResponse
// @Router       /webhooks/handleEvents [post]
// @BasePath /
func (h *handler) HandleClerkEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, ok := ctx.Value(webhookEventKey).(types.WebhookEvent)
	if !ok {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: "No webhook event in context",
		}, h.appCtx.Logger)
		return
	}

	h.appCtx.Logger.Info("Handling webhook event", "event", event)

	err := h.webhookService.HandleEvent(ctx, &event)

	if err != nil {
		response.WriteJSONError(w, err, h.appCtx.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
}

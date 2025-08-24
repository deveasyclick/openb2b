package user

import (
	"encoding/json"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type Handler struct {
	service interfaces.UserService
	appCtx  *deps.AppContext
}

func NewHandler(service interfaces.UserService, appCtx *deps.AppContext) interfaces.UserHandler {
	return &Handler{service: service, appCtx: appCtx}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userFromContext, err := identity.UserFromContext(ctx)

	if err != nil {
		h.appCtx.Logger.Error("error getting user from context", "err", err)
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, h.appCtx.Logger)
		return
	}

	user, apiError := h.service.FindByID(ctx, userFromContext.ID, []string{"Org"})
	if apiError != nil {
		response.WriteJSONError(w, apiError, h.appCtx.Logger)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.appCtx.Logger.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

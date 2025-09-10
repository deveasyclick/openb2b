package user

import (
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type APIResponseUser struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    model.User `json:"data"`
}

type Handler struct {
	service interfaces.UserService
	appCtx  *deps.AppContext
}

func NewHandler(service interfaces.UserService, appCtx *deps.AppContext) interfaces.UserHandler {
	return &Handler{service: service, appCtx: appCtx}
}

// Get godoc
// @Summary Get authenticated user
// @Description Get an authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} APIResponseUser
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 401  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /users/me [get]
// @Security BearerAuth
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userFromContext, err := identity.UserFromContext(ctx)

	if err != nil {
		response.WriteJSONErrorV2(w, 401, nil, "unauthorized", h.appCtx.Logger)
		return
	}

	user, err := h.service.FindByID(ctx, userFromContext.ID, []string{"Org"})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.WriteJSONErrorV2(w, 404, err, "user not found", h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, 500, err, apperrors.ErrFindUser, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, user, h.appCtx.Logger)
}

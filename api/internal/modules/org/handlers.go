package org

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/internal/shared/validator"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

type OrgHandler struct {
	service     interfaces.OrgService
	createOrgUC interfaces.CreateOrgUseCase
	appCtx      *deps.AppContext
}

func NewHandler(service interfaces.OrgService, createOrgUC interfaces.CreateOrgUseCase, appCtx *deps.AppContext) interfaces.OrgHandler {
	return &OrgHandler{service: service, createOrgUC: createOrgUC, appCtx: appCtx}
}

// Create godoc
// @Summary Create organization
// @Description Create a new organization
// @Tags organizations
// @Accept json
// @Produce json
// @Param request body createDTO true "Organization payload"
// @Success 200 {object} model.Org
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /orgs [post]
// @Security BearerAuth
func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		h.appCtx.Logger.Error("invalid request body in org create", "errors", errors)
		validator.WriteValidationResponse(w, errors)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		h.appCtx.Logger.Error(apperrors.ErrUserFromContext, "err", err)
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateOrg,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrUserFromContext, err),
		}, h.appCtx.Logger)
		return
	}

	// Convert request to model
	org := req.ToModel()

	// Check if org already exists
	exists, apiError := h.service.Exists(ctx, map[string]any{"name": org.Name})
	if apiError != nil {
		response.WriteJSONError(w, apiError, h.appCtx.Logger)
		return
	}

	// Return already exists error if org already exists
	if exists {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("%s: name %s", apperrors.ErrOrgAlreadyExists, org.Name),
		}, h.appCtx.Logger)
		return
	}

	apiError = h.createOrgUC.Execute(ctx, types.CreateOrgInput{
		Org:  org,
		User: userFromContext,
	})

	if apiError != nil {
		response.WriteJSONError(w, apiError, h.appCtx.Logger)
		return
	}

	if err := json.NewEncoder(w).Encode(org); err != nil {
		h.appCtx.Logger.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

// Update godoc
// @Summary Update organization
// @Description Update an existing organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param request body updateDTO true "Update organization payload"
// @Success 200 {object} model.Org
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orgs/{id} [put]
// @Security BearerAuth
func (h *OrgHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %d", apperrors.ErrInvalidId, id), http.StatusBadRequest)
		return
	}

	var req updateDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing org
	existingOrg, apiError := h.service.FindOrg(ctx, uint(id))
	if apiError != nil {
		response.WriteJSONError(w, apiError, h.appCtx.Logger)
		return
	}

	// Update only provided fields
	req.ApplyModel(existingOrg)

	if apiError := h.service.Update(ctx, existingOrg); apiError != nil {
		response.WriteJSONError(w, apiError, h.appCtx.Logger)
		return
	}

	if err := json.NewEncoder(w).Encode(existingOrg); err != nil {
		h.appCtx.Logger.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

// Delete godoc
// @Summary Delete organization
// @Description Delete an organization by ID
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {integer} int "Deleted organization ID"
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orgs/{id} [delete]
// @Security BearerAuth
func (h *OrgHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidId, http.StatusBadRequest)
		return
	}

	if apiError := h.service.Delete(ctx, uint(id)); apiError != nil {
		response.WriteJSONError(w, apiError, h.appCtx.Logger)
		return
	}

	if err := json.NewEncoder(w).Encode(id); err != nil {
		h.appCtx.Logger.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

// Get godoc
// @Summary Get organization
// @Description Get an organization by ID
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} model.Org
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orgs/{id} [get]
// @Security BearerAuth
func (h *OrgHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidId, http.StatusBadRequest)
		return
	}

	org, apiError := h.service.FindOrg(ctx, uint(id))
	if apiError != nil {
		response.WriteJSONError(w, apiError, h.appCtx.Logger)
		return
	}

	if err := json.NewEncoder(w).Encode(org); err != nil {
		h.appCtx.Logger.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

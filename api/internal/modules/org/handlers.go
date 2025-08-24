package org

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
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

func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		h.appCtx.Logger.Error("invalid request body in org create", "errors", errors)
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Convert request to model
	org := req.ToModel()

	err := h.createOrgUC.Execute(ctx, types.CreateOrgInput{
		Org:    org,
		UserID: 1,
	})

	if err != nil {
		response.WriteJSONError(w, err, h.appCtx.Logger)
		return
	}

	if err := json.NewEncoder(w).Encode(org); err != nil {
		h.appCtx.Logger.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

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

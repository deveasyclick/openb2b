package org

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/shared/httphelper"
	"github.com/deveasyclick/openb2b/internal/shared/validator"
	"github.com/deveasyclick/openb2b/pkg/apperrors"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

type orgHandler struct {
	service     interfaces.OrgService
	createOrgUC CreateOrgUseCase
}

func NewOrgHandler(service interfaces.OrgService, createOrgUC CreateOrgUseCase) interfaces.OrgHandler {
	return &orgHandler{service: service, createOrgUC: createOrgUC}
}

func (h *orgHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateOrgDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		slog.Error("invalid request body in org create", "errors", errors)
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Convert request to model
	org := req.ToModel()

	err := h.createOrgUC.Execute(ctx, CreateOrgInput{
		Org:    org,
		UserID: 1,
	})

	if err != nil {
		httphelper.WriteJSONError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(org); err != nil {
		slog.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

func (h *orgHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %d", apperrors.ErrInvalidId, id), http.StatusBadRequest)
		return
	}

	var req UpdateOrgDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing org
	existingOrg, apiError := h.service.FindOrg(ctx, uint(id))
	if err != nil {
		httphelper.WriteJSONError(w, apiError)
		return
	}

	// Update only provided fields
	req.ApplyModel(existingOrg)

	if apiError := h.service.Update(ctx, existingOrg); apiError != nil {
		httphelper.WriteJSONError(w, apiError)
		return
	}

	if err := json.NewEncoder(w).Encode(existingOrg); err != nil {
		slog.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

func (h *orgHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidId, http.StatusBadRequest)
		return
	}

	if apiError := h.service.Delete(ctx, uint(id)); apiError != nil {
		httphelper.WriteJSONError(w, apiError)
		return
	}

	if err := json.NewEncoder(w).Encode(id); err != nil {
		slog.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

func (h *orgHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidId, http.StatusBadRequest)
		return
	}

	org, apiError := h.service.FindOrg(ctx, uint(id))
	if apiError != nil {
		httphelper.WriteJSONError(w, apiError)
		return
	}

	if err := json.NewEncoder(w).Encode(org); err != nil {
		slog.Warn(apperrors.ErrEncodeResponse, "error", err)
	}
}

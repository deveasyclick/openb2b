package org

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/internal/shared/validator"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type APIResponseOrg struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    model.Org `json:"data"`
}

type APIResponseInt struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    int    `json:"data"`
}

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
// @Param request body dto.CreateOrgDTO true "Organization payload"
// @Success 200 {object} APIResponseOrg
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /orgs [post]
// @Security BearerAuth
func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateOrgDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateOrg, h.appCtx.Logger)
		return
	}

	// Convert request to model
	org := req.ToModel()

	// Check if org already exists
	exists, err := h.service.Exists(ctx, map[string]any{"name": org.Name})
	if err != nil && err != gorm.ErrRecordNotFound {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrOrgAlreadyExists, h.appCtx.Logger)
		return
	}

	// Return already exists error if org already exists
	if exists {
		response.WriteJSONErrorV2(w, http.StatusConflict, nil, apperrors.ErrOrgAlreadyExists, h.appCtx.Logger)
		return
	}

	err = h.createOrgUC.Execute(ctx, types.CreateOrgInput{
		Org:  org,
		User: userFromContext,
	})

	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateOrg, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusCreated, org, h.appCtx.Logger)
}

// Update godoc
// @Summary Update organization
// @Description Update an existing organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param request body dto.UpdateOrgDTO true "Update organization payload"
// @Success 200 {object} APIResponseOrg
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orgs/{id} [patch]
// @Security BearerAuth
func (h *OrgHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	var req dto.UpdateOrgDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing org
	existingOrg, err := h.service.FindOrg(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrOrgNotFound, h.appCtx.Logger)
			return
		}
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateOrg, h.appCtx.Logger)

		return
	}

	// Update only provided fields
	req.ApplyModel(existingOrg)

	if err := h.service.Update(ctx, existingOrg); err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateOrg, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, existingOrg, h.appCtx.Logger)
}

// Delete godoc
// @Summary Delete organization
// @Description Delete an organization by ID
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {integer} APIResponseInt
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orgs/{id} [delete]
// @Security BearerAuth
func (h *OrgHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}
	if err := h.service.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrOrgNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrDeleteOrg, h.appCtx.Logger)
		return
	}
	response.WriteJSONSuccess(w, http.StatusOK, id, h.appCtx.Logger)
}

// Get godoc
// @Summary Get organization
// @Description Get an organization by ID
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} APIResponseOrg
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orgs/{id} [get]
// @Security BearerAuth
func (h *OrgHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	org, err := h.service.FindOrg(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrOrgNotFound, h.appCtx.Logger)
			return
		}
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFindOrg, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, org, h.appCtx.Logger)
}

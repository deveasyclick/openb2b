package customer

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/validator"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

var allowedCustomerSearchFields = map[string]bool{"first_name": true, "last_name": true, "phone_number": true, "email": true, "company": true}

// For Swagger docs
type APIResponseCustomer struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    model.Customer `json:"data"`
}

type CustomerHandler struct {
	service interfaces.CustomerService
	appCtx  *deps.AppContext
}

func NewHandler(service interfaces.CustomerService, appCtx *deps.AppContext) interfaces.CustomerHandler {
	return &CustomerHandler{service: service, appCtx: appCtx}
}

// Filter godoc
// @Summary      List customers with filtering and pagination
// @Description  Returns a paginated list of customers. Supports filtering, sorting, searching, and preloading.
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        page          query     int     false  "Page number (default: 1)"
// @Param        limit         query     int     false  "Number of items per page (default: 20, max: 100)"
// @Param        sort          query     string  false  "Sort by field, e.g. 'created_at desc'"
// @Param        preloads      query     string  false  "Comma-separated list of relations to preload. relation must start with uppercase. e.g. 'Orders,Org'"
// @Param        search_fields query     string  false  "Comma-separated list of fields to search (must be allowed)"
// @Param        first_name          query     string  false  "Filter by first name"
// @Param        last_name     query     string  false  "Filter by last name"
// @Param        phone_number  query     string  false  "Filter by phone number"
// @Param        email         query     string  false  "Filter by email"
// @Param        company       query     string  false  "Filter by company"
// @Success      200           {object}  APIResponseCustomer
// @Failure      400           {object}  apperrors.APIError "Invalid filter parameters"
// @Failure      500           {object}  apperrors.APIError "Internal server error"
// @Router       /customers [get]
// @Security BearerAuth
func (h *CustomerHandler) Filter(w http.ResponseWriter, r *http.Request) {
	opts, err := pagination.ParsePaginationOptions(r.URL.Query(), allowedCustomerSearchFields)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrFilterCustomer, h.appCtx.Logger)
		return
	}

	customers, total, err := h.service.Filter(r.Context(), opts)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFilterCustomer, h.appCtx.Logger)
		return
	}

	resp := response.FilterResponse[model.Customer]{
		Pagination: pagination.BuildPagination(total, opts),
		Items:      customers,
	}

	response.WriteJSONSuccess(w, http.StatusOK, resp, h.appCtx.Logger)
}

// Create godoc
// @Summary Create customers
// @Description Create a new customer
// @Tags customers
// @Accept json
// @Produce json
// @Param request body dto.CreateCustomerDTO true "Customer payload"
// @Success 200 {object} APIResponseCustomer
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      409  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /customers [post]
// @Security BearerAuth
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateCustomerDTO
	if errs := validator.ValidateRequest(r, &req); len(errs) > 0 {
		validator.WriteValidationResponse(w, errs)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateCustomer, h.appCtx.Logger)
		return
	}

	customer := req.ToModel(userFromContext.Org)
	err = h.service.Create(ctx, customer)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateCustomer, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusCreated, customer, h.appCtx.Logger)
}

// Update godoc
// @Summary Update customer
// @Description Update an existing customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param request body dto.UpdateCustomerDTO true "Update customer payload"
// @Success 200 {object} APIResponseCustomer
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /customers/{id} [patch]
// @Security BearerAuth
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	var req dto.UpdateCustomerDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	existingCustomer, err := h.service.Update(ctx, uint(id), &req)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrCustomerNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateCustomer, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, existingCustomer, h.appCtx.Logger)
}

// Delete godoc
// @Summary Delete customer
// @Description Delete a customer by ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {integer} response.APIResponseInt
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /customers/{id} [delete]
// @Security BearerAuth
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrCustomerNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrDeleteCustomer, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, id, h.appCtx.Logger)
}

// Get godoc
// @Summary Get customer
// @Description Get a customer by ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} APIResponseCustomer
// @Failure 400 {object} apperrors.APIErrorResponse
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 500 {object} apperrors.APIErrorResponse
// @Router /customers/{id} [get]
// @Security BearerAuth
func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	customer, err := h.service.FindOneWithFields(ctx, nil, map[string]any{"id": id}, []string{"Org", "Orders"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrCustomerNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFindCustomer, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, customer, h.appCtx.Logger)
}

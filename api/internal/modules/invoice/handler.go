package invoice

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

var allowedInvoiceSearchFields = map[string]bool{"invoice_number": true, "notes": true, "status": true}

// For Swagger docs
type APIResponseInvoice struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    model.Invoice `json:"data"`
}

type InvoiceHandler struct {
	service interfaces.InvoiceService
	appCtx  *deps.AppContext
}

func NewHandler(service interfaces.InvoiceService, appCtx *deps.AppContext) interfaces.InvoiceHandler {
	return &InvoiceHandler{service: service, appCtx: appCtx}
}

// Filter godoc
// @Summary      List invoices with filtering and pagination
// @Description  Returns a paginated list of invoices. Supports filtering, sorting, searching, and preloading.
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        page          query     int     false  "Page number (default: 1)"
// @Param        limit         query     int     false  "Number of items per page (default: 20, max: 100)"
// @Param        sort          query     string  false  "Sort by field, e.g. 'created_at desc'"
// @Param        preloads      query     string  false  "Comma-separated list of relations to preload. relation must start with uppercase. e.g. 'Orders,Org'"
// @Param        search_fields query     string  false  "Comma-separated list of fields to search (must be allowed)"
// @Param        order_number          query     string  false  "Filter by order number"
// @Param        notes     query     string  false  "Filter by notes"
// @Param        status       query     string  false  "Filter by status"
// @Success      200           {object}  APIResponseInvoice
// @Failure      400           {object}  apperrors.APIError "Invalid filter parameters"
// @Failure      500           {object}  apperrors.APIError "Internal server error"
// @Router       /invoices [get]
// @Security BearerAuth
func (h *InvoiceHandler) Filter(w http.ResponseWriter, r *http.Request) {
	opts, err := pagination.ParsePaginationOptions(r.URL.Query(), allowedInvoiceSearchFields)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrFilterInvoice, h.appCtx.Logger)
		return
	}

	invoices, total, err := h.service.Filter(r.Context(), opts)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFilterInvoice, h.appCtx.Logger)
		return
	}

	resp := response.FilterResponse[model.Invoice]{
		Pagination: pagination.BuildPagination(total, opts),
		Items:      invoices,
	}

	response.WriteJSONSuccess(w, http.StatusOK, resp, h.appCtx.Logger)
}

// Create godoc
// @Summary Create invoices
// @Description Create a new invoice
// @Tags invoices
// @Accept json
// @Produce json
// @Param request body dto.CreateInvoiceDTO true "Invoice payload"
// @Success 200 {object} APIResponseInvoice
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      409  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /invoices [post]
// @Security BearerAuth
func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateInvoiceDTO
	if errs := validator.ValidateRequest(r, &req); len(errs) > 0 {
		validator.WriteValidationResponse(w, errs)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateInvoice, h.appCtx.Logger)
		return
	}

	invoice, err := h.service.Create(ctx, userFromContext.Org, &req)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateInvoice, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusCreated, invoice, h.appCtx.Logger)
}

// Update godoc
// @Summary Update invoice
// @Description Update an existing invoice by ID
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param request body dto.UpdateInvoiceDTO true "Update invoice payload"
// @Success 200 {object} APIResponseInvoice
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /invoices/{id} [patch]
// @Security BearerAuth
func (h *InvoiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	var req dto.UpdateInvoiceDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	existingInvoice, err := h.service.Update(ctx, uint(id), &req)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrInvoiceNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateInvoice, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, existingInvoice, h.appCtx.Logger)
}

// Delete godoc
// @Summary Delete invoice
// @Description Delete a invoice by ID
// @Tags invoices
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {integer} response.APIResponseInt
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /invoices/{id} [delete]
// @Security BearerAuth
func (h *InvoiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrInvoiceNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrDeleteInvoice, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, id, h.appCtx.Logger)
}

// Get godoc
// @Summary Get invoice
// @Description Get a invoice by ID
// @Tags invoices
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} APIResponseInvoice
// @Failure 400 {object} apperrors.APIErrorResponse
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 500 {object} apperrors.APIErrorResponse
// @Router /invoices/{id} [get]
// @Security BearerAuth
func (h *InvoiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	invoice, err := h.service.FindOneWithFields(ctx, nil, map[string]any{"id": id}, []string{"Items", "Order"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrInvoiceNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFindInvoice, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, invoice, h.appCtx.Logger)
}

// Update godoc
// @Summary Update invoice
// @Description Issue an invoice by ID
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} response.APIResponseString
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /invoices/{id}/issue [post]
// @Security BearerAuth
func (h *InvoiceHandler) Issue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	err = h.service.Issue(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrInvoiceNotFound, h.appCtx.Logger)
			return
		}

		if errors.Is(err, errors.New(apperrors.ErrInvalidInvoiceStatus)) {
			response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidInvoiceStatus, h.appCtx.Logger)
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrIssueInvoice, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, "Invoice issued and emailed", h.appCtx.Logger)
}

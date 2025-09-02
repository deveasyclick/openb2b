package order

import (
	"errors"
	"fmt"
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

var allowedOrderSearchFields = map[string]bool{"order_number": true, "notes": true}

// For Swagger docs
type APIResponseOrder struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    model.Order `json:"data"`
}

type OrderHandler struct {
	service interfaces.OrderService
	appCtx  *deps.AppContext
}

func NewHandler(service interfaces.OrderService, appCtx *deps.AppContext) interfaces.OrderHandler {
	return &OrderHandler{service: service, appCtx: appCtx}
}

// Filter godoc
// @Summary      List orders with filtering and pagination
// @Description  Returns a paginated list of orders. Supports filtering, sorting, searching, and preloading.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        page          query     int     false  "Page number (default: 1)"
// @Param        limit         query     int     false  "Number of items per page (default: 20, max: 100)"
// @Param        sort          query     string  false  "Sort by field, e.g. 'created_at desc'"
// @Param        preloads      query     string  false  "Comma-separated list of relations to preload. relation must start with uppercase"
// @Param        search_fields query     string  false  "Comma-separated list of fields to search (must be allowed)"
// @Param        order_number          query     string  false  "Filter by order number"
// @Param        notes     query     string  false  "Filter by notes"
// @Success      200           {object}  APIResponseOrder
// @Failure      400           {object}  apperrors.APIError "Invalid filter parameters"
// @Failure      500           {object}  apperrors.APIError "Internal server error"
// @Router       /orders [get]
// @Security BearerAuth
func (h *OrderHandler) Filter(w http.ResponseWriter, r *http.Request) {
	opts, err := pagination.ParsePaginationOptions(r.URL.Query(), allowedOrderSearchFields)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrFilterOrder, h.appCtx.Logger)
		return
	}

	orders, total, err := h.service.Filter(r.Context(), opts)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFilterOrder, h.appCtx.Logger)
		return
	}

	resp := response.FilterResponse[model.Order]{
		Pagination: pagination.BuildPagination(total, opts),
		Items:      orders,
	}

	response.WriteJSONSuccess(w, http.StatusOK, resp, h.appCtx.Logger)
}

// Create godoc
// @Summary Create orders
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderDTO true "Order payload"
// @Success 200 {object} APIResponseOrder
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      409  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /orders [post]
// @Security BearerAuth
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateOrderDTO
	if errs := validator.ValidateRequest(r, &req); len(errs) > 0 {
		validator.WriteValidationResponse(w, errs)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateOrder, h.appCtx.Logger)
		return
	}

	order, err := h.service.Create(ctx, req, userFromContext.Org)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateOrder, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusCreated, order, h.appCtx.Logger)
}

// Update godoc
// @Summary Update order
// @Description Update an existing order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body dto.UpdateOrderDTO true "Update order payload"
// @Success 200 {object} APIResponseOrder
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orders/{id} [patch]
// @Security BearerAuth
func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	var req dto.UpdateOrderDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing order
	existingOrder, err := h.service.FindByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrOrderNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateOrder, h.appCtx.Logger)
		return
	}

	// Don't update order if it is not in pending status
	if existingOrder.Status != model.OrderStatusPending {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, fmt.Sprintf("cannot update %s: %s", apperrors.ErrOrderNotPending, existingOrder.Status), h.appCtx.Logger)

		return
	}

	if err := h.service.Update(ctx, existingOrder, req); err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateOrder, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, existingOrder, h.appCtx.Logger)
}

// Delete godoc
// @Summary Delete order
// @Description Delete a order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {integer} response.APIResponseInt
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /orders/{id} [delete]
// @Security BearerAuth
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrOrderNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrDeleteOrder, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, id, h.appCtx.Logger)
}

// Get godoc
// @Summary Get order
// @Description Get a order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} APIResponseOrder
// @Failure 400 {object} apperrors.APIErrorResponse
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 500 {object} apperrors.APIErrorResponse
// @Router /orders/{id} [get]
// @Security BearerAuth
func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	order, err := h.service.FindOneWithFields(ctx, nil, map[string]any{"id": id}, []string{"Items", "Customer"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrOrderNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFindOrder, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, order, h.appCtx.Logger)
}

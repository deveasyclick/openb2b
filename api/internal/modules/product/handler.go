package product

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

var allowedProductSearchFields = map[string]bool{"name": true, "last_name": true, "phone_number": true, "email": true}

// For Swagger docs
type APIResponseProduct struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    model.Product `json:"data"`
}

type APIResponseVariant struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    model.Variant `json:"data"`
}
type ProductHandler struct {
	service interfaces.ProductService
	appCtx  *deps.AppContext
}

func NewHandler(service interfaces.ProductService, appCtx *deps.AppContext) interfaces.ProductHandler {
	return &ProductHandler{service: service, appCtx: appCtx}
}

// Filter godoc
// @Summary      List products with filtering and pagination
// @Description  Returns a paginated list of products. Supports filtering, sorting, searching, and preloading.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page          query     int     false  "Page number (default: 1)"
// @Param        limit         query     int     false  "Number of items per page (default: 20, max: 100)"
// @Param        sort          query     string  false  "Sort by field, e.g. 'created_at desc'"
// @Param        preloads      query     string  false  "Comma-separated list of relations to preload. relation must start with uppercase"
// @Param        search_fields query     string  false  "Comma-separated list of fields to search (must be allowed)"
// @Param        name          query     string  false  "Filter by product name"
// @Param        last_name     query     string  false  "Filter by last name"
// @Param        phone_number  query     string  false  "Filter by phone number"
// @Param        email         query     string  false  "Filter by email"
// @Success      200           {object}  APIResponseProduct
// @Failure      400           {object}  apperrors.APIError "Invalid filter parameters"
// @Failure      500           {object}  apperrors.APIError "Internal server error"
// @Router       /products [get]
// @Security BearerAuth
func (h *ProductHandler) Filter(w http.ResponseWriter, r *http.Request) {
	opts, err := pagination.ParsePaginationOptions(r.URL.Query(), allowedProductSearchFields)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrFilterProduct, h.appCtx.Logger)
		return
	}

	products, total, err := h.service.Filter(r.Context(), opts)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFilterProduct, h.appCtx.Logger)
		return
	}

	resp := response.FilterResponse[model.Product]{
		Pagination: pagination.BuildPagination(total, opts),
		Items:      products,
	}

	response.WriteJSONSuccess(w, http.StatusOK, resp, h.appCtx.Logger)
}

// Create godoc
// @Summary Create products
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductDTO true "Product payload"
// @Success 200 {object} APIResponseProduct
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      409  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /products [post]
// @Security BearerAuth
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateProductDTO
	if errs := validator.ValidateRequest(r, &req); len(errs) > 0 {
		validator.WriteValidationResponse(w, errs)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateProduct, h.appCtx.Logger)
		return
	}

	// Convert request to model
	product := req.ToModel(userFromContext.Org)

	// Check if product already exists
	exists, err := h.service.Exists(ctx, map[string]any{"name": product.Name})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateProduct, h.appCtx.Logger)
		return
	}

	// Return already exists error if product already exists
	if exists {
		response.WriteJSONErrorV2(w, http.StatusConflict, nil, fmt.Sprintf("%s: name %s", apperrors.ErrProductAlreadyExists, product.Name), h.appCtx.Logger)
		return
	}

	if err = h.service.Create(ctx, &product); err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateProduct, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusCreated, product, h.appCtx.Logger)
}

// Update godoc
// @Summary Update product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductDTO true "Update product payload"
// @Success 200 {object} APIResponseProduct
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id} [patch]
// @Security BearerAuth
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	var req dto.UpdateProductDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing product
	existingProduct, err := h.service.FindByID(ctx, uint(id))
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateProduct, h.appCtx.Logger)
		return
	}

	// Update only provided fields
	req.ApplyModel(existingProduct)

	if err := h.service.Update(ctx, existingProduct); err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateProduct, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, existingProduct, h.appCtx.Logger)
}

// Delete godoc
// @Summary Delete product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {integer} response.APIResponseInt
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrProductNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrDeleteProduct, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, id, h.appCtx.Logger)
}

// Get godoc
// @Summary Get product
// @Description Get a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} APIResponseProduct
// @Failure 400 {object} apperrors.APIErrorResponse
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 500 {object} apperrors.APIErrorResponse
// @Router /products/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	product, err := h.service.FindOneWithFields(ctx, nil, map[string]any{"id": id}, []string{"Variants"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrProductNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFindProduct, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, product, h.appCtx.Logger)
}

// ------------------------Variants-----------------------
// Create godoc
// @Summary Create variant
// @Description Create a new variant
// @Tags variants
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Param request body dto.CreateProductVariantDTO true "Variant payload"
// @Success 200 {object} APIResponseVariant
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      409  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /products/{productId}/variants [post]
// @Security BearerAuth
func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseUint(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, apperrors.ErrInvalidId, h.appCtx.Logger)
		return
	}

	ctx := r.Context()
	var req dto.CreateProductVariantDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateVariant, h.appCtx.Logger)
		return
	}

	// Convert request to model
	variant := req.ToModel(userFromContext.Org)
	variant.ProductID = uint(productId)

	// Check if variant already exists
	exists, err := h.service.CheckVariantExists(ctx, variant.SKU)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateVariant, h.appCtx.Logger)
		return
	}

	// Return already exists error if variant already exists
	if exists {
		response.WriteJSONErrorV2(w, http.StatusConflict, nil, fmt.Sprintf("%s: sku %s", apperrors.ErrVariantAlreadyExists, variant.SKU), h.appCtx.Logger)
		return
	}

	if err = h.service.CreateVariant(ctx, &variant); err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrCreateVariant, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusCreated, variant, h.appCtx.Logger)
}

// Update godoc
// @Summary Update variant
// @Description Update an existing variant by product ID and variant ID
// @Tags variants
// @Accept json
// @Produce json
// @param productId path int true "Product ID"
// @Param id path int true "Variant ID"
// @Param request body dto.UpdateVariantDTO true "Update variant payload"
// @Success 200 {object} APIResponseVariant
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{productId}/variants/{id} [patch]
// @Security BearerAuth
func (h *ProductHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productId, err := strconv.ParseUint(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, fmt.Sprintf("%s: productId %d", apperrors.ErrInvalidId, productId), h.appCtx.Logger)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, fmt.Sprintf("%s: variantId %d", apperrors.ErrInvalidId, productId), h.appCtx.Logger)
		return
	}

	var req dto.UpdateVariantDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing variant
	existingVariant, err := h.service.FindVariantByID(ctx, uint(productId), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrVariantNotFound, h.appCtx.Logger)
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateVariant, h.appCtx.Logger)
		return
	}

	// Update only provided fields
	req.ApplyModel(existingVariant)

	if err := h.service.UpdateVariant(ctx, existingVariant); err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrUpdateVariant, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, existingVariant, h.appCtx.Logger)
}

// Delete godoc
// @Summary Delete variant
// @Description Delete a variant by product ID and variant ID
// @Tags variants
// @Produce json
// @Param productId path int true "Product ID"
// @Param id path int true "Variant ID"
// @Success 200 {integer} response.APIResponseInt
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id}/variants/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productId, err := strconv.ParseUint(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, fmt.Sprintf("%s: productId %d", apperrors.ErrInvalidId, productId), h.appCtx.Logger)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, fmt.Sprintf("%s: variantId %d", apperrors.ErrInvalidId, productId), h.appCtx.Logger)
		return
	}

	if err := h.service.DeleteVariant(ctx, uint(productId), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.WriteJSONErrorV2(w, http.StatusNotFound, nil, apperrors.ErrVariantNotFound, h.appCtx.Logger)
			return
		}

		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrDeleteVariant, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, id, h.appCtx.Logger)
}

// Get godoc
// @Summary Get variant
// @Description Get a variant by product ID and variant ID
// @Tags variants
// @Produce json
// @Param productId path int true "Product ID"
// @Param id path int true "variant ID"
// @Success 200 {object} APIResponseVariant
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id}/variants/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) GetVariant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productId, err := strconv.ParseUint(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, fmt.Sprintf("%s: productId %d", apperrors.ErrInvalidId, productId), h.appCtx.Logger)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusBadRequest, nil, fmt.Sprintf("%s: variantId %d", apperrors.ErrInvalidId, productId), h.appCtx.Logger)
		return
	}

	variant, err := h.service.FindVariantByID(ctx, uint(productId), uint(id))
	if err != nil {
		response.WriteJSONErrorV2(w, http.StatusInternalServerError, err, apperrors.ErrFindVariant, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, variant, h.appCtx.Logger)
}

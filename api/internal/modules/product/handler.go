package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/validator"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	service interfaces.ProductService
	appCtx  *deps.AppContext
}

func NewHandler(service interfaces.ProductService, appCtx *deps.AppContext) interfaces.ProductHandler {
	return &ProductHandler{service: service, appCtx: appCtx}
}

// Create godoc
// @Summary Create products
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body createProductDTO true "Product payload"
// @Success 200 {object} model.Product
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /products [post]
// @Security BearerAuth
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createProductDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		h.appCtx.Logger.Error(apperrors.ErrInvalidRequestBody, "errors", errors)
		validator.WriteValidationResponse(w, errors)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		h.appCtx.Logger.Error(apperrors.ErrUserFromContext, "err", err)
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateProduct,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrUserFromContext, err),
		}, h.appCtx.Logger)
		return
	}

	// Convert request to model
	product := req.ToModel()
	product.OrgID = userFromContext.Org

	// Check if product already exists
	exists, err := h.service.Exists(ctx, map[string]any{"name": product.Name})
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateProduct,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrCreateProduct, err),
		}, h.appCtx.Logger)
		return
	}

	// Return already exists error if product already exists
	if exists {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("%s: name %s", apperrors.ErrProductAlreadyExists, product.Name),
		}, h.appCtx.Logger)
		return
	}

	if err = h.service.Create(ctx, &product); err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateProduct,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrCreateProduct, err),
		}, h.appCtx.Logger)
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
// @Param request body updateProductDTO true "Update product payload"
// @Success 200 {object} model.Product
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("%s: %d", apperrors.ErrInvalidId, id),
		}, h.appCtx.Logger)
		return
	}

	var req updateProductDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing product
	existingProduct, err := h.service.FindByID(ctx, uint(id))
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrUpdateProduct,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrUpdateProduct, err),
		}, h.appCtx.Logger)
		return
	}

	// Update only provided fields
	req.ApplyModel(existingProduct)

	if err := h.service.Update(ctx, existingProduct); err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrUpdateProduct,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrUpdateProduct, err),
		}, h.appCtx.Logger)
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
// @Success 200 {integer} int "Deleted product ID"
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: apperrors.ErrInvalidId,
		}, h.appCtx.Logger)
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrDeleteProduct,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrDeleteProduct, err),
		}, h.appCtx.Logger)
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
// @Success 200 {object} model.Product
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: apperrors.ErrInvalidId,
		}, h.appCtx.Logger)
		return
	}

	product, err := h.service.FindByID(ctx, uint(id))
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrFindProduct,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrFindProduct, err),
		}, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, product, h.appCtx.Logger)
}

// ------------------------Variants-----------------------
// Create godoc
// @Summary Create variant
// @Description Create a new variant
// @Tags products
// @Accept json
// @Produce json
// @Param request body createVariantDTO true "Variant payload"
// @Success 200 {object} model.Variant
// @Failure      400  {object}  apperrors.APIErrorResponse
// @Failure      500  {object}  apperrors.APIErrorResponse
// @Router /products/{id}/variants [post]
// @Security BearerAuth
func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createVariantDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		h.appCtx.Logger.Error(apperrors.ErrInvalidRequestBody, "errors", errors)
		validator.WriteValidationResponse(w, errors)
		return
	}

	userFromContext, err := identity.UserFromContext(ctx)
	if err != nil {
		h.appCtx.Logger.Error(apperrors.ErrUserFromContext, "err", err)
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateProduct,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrUserFromContext, err),
		}, h.appCtx.Logger)
		return
	}

	// Convert request to model
	variant := req.ToModel()
	variant.OrgID = userFromContext.Org

	// Check if variant already exists
	exists, err := h.service.Exists(ctx, map[string]any{"sku": variant.SKU})
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateVariant,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrCreateVariant, err),
		}, h.appCtx.Logger)
		return
	}

	// Return already exists error if variant already exists
	if exists {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("%s: sku %s", apperrors.ErrProductAlreadyExists, variant.SKU),
		}, h.appCtx.Logger)
		return
	}

	if err = h.service.CreateVariant(ctx, &variant); err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateVariant,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrCreateVariant, err),
		}, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusCreated, variant, h.appCtx.Logger)
}

// Update godoc
// @Summary Update variant
// @Description Update an existing variant by product ID and variant ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Variant ID"
// @params productId path int true "Product ID"
// @Param request body updateVariantDTO true "Update variant payload"
// @Success 200 {object} model.Variant
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{productId}/variants/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productId, err := strconv.ParseUint(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("%s: %d", apperrors.ErrInvalidId, productId),
		}, h.appCtx.Logger)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("%s: %d", apperrors.ErrInvalidId, id),
		}, h.appCtx.Logger)
		return
	}

	var req updateVariantDTO
	if errors := validator.ValidateRequest(r, &req); len(errors) > 0 {
		validator.WriteValidationResponse(w, errors)
		return
	}

	// Get existing variant
	existingVariant, err := h.service.FindVariantByID(ctx, uint(productId), uint(id))
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrUpdateVariant,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrUpdateVariant, err),
		}, h.appCtx.Logger)
		return
	}

	// Update only provided fields
	req.ApplyModel(existingVariant)

	if err := h.service.UpdateVariant(ctx, existingVariant); err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrUpdateVariant,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrUpdateVariant, err),
		}, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, existingVariant, h.appCtx.Logger)
}

// Delete godoc
// @Summary Delete variant
// @Description Delete a variant by product ID and variant ID
// @Tags products
// @Produce json
// @Param productId path int true "Product ID"
// @Param id path int true "Variant ID"
// @Success 200 {integer} int "Deleted variant ID"
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id}/variants/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productId, err := strconv.ParseUint(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: apperrors.ErrInvalidId,
		}, h.appCtx.Logger)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: apperrors.ErrInvalidId,
		}, h.appCtx.Logger)
		return
	}

	if err := h.service.DeleteVariant(ctx, uint(productId), uint(id)); err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrDeleteVariant,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrDeleteVariant, err),
		}, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, id, h.appCtx.Logger)
}

// Get godoc
// @Summary Get variant
// @Description Get a variant by product ID and variant ID
// @Tags products
// @Produce json
// @Param productId path int true "Product ID"
// @Param id path int true "variant ID"
// @Success 200 {object} model.Variant
// @Failure 404 {object} apperrors.APIErrorResponse
// @Failure 400  {object}  apperrors.APIErrorResponse
// @Failure 500  {object}  apperrors.APIErrorResponse
// @Router /products/{id}/variants/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) GetVariant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productId, err := strconv.ParseUint(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: apperrors.ErrInvalidId,
		}, h.appCtx.Logger)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: apperrors.ErrInvalidId,
		}, h.appCtx.Logger)
		return
	}

	variant, err := h.service.FindVariantByID(ctx, uint(productId), uint(id))
	if err != nil {
		response.WriteJSONError(w, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrFindVariant,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrFindVariant, err),
		}, h.appCtx.Logger)
		return
	}

	response.WriteJSONSuccess(w, http.StatusOK, variant, h.appCtx.Logger)
}

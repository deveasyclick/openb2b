package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"gorm.io/gorm"
)

type ProductService interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, ID uint) error
	FindByID(ctx context.Context, ID uint) (*model.Product, error)
	Exists(ctx context.Context, where map[string]any) (bool, error)
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Product, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Product, int64, error)
	WithTx(tx *gorm.DB) ProductService

	// Varaiants
	CreateVariant(ctx context.Context, variant *model.Variant) error
	UpdateVariant(ctx context.Context, variant *model.Variant) error
	DeleteVariant(ctx context.Context, productID uint, variantID uint) error
	FindVariantByID(ctx context.Context, productID uint, variantID uint) (*model.Variant, error)
	CheckVariantExists(ctx context.Context, sku string) (bool, error)
	FindVariants(ctx context.Context, where map[string]any, preloads []string) ([]model.Variant, error)
}

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, ID uint) (*model.Product, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Product, int64, error)
	Delete(ctx context.Context, ID uint) error
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Product, error)
	WithTx(tx *gorm.DB) ProductRepository

	// Variants
	CreateVariant(ctx context.Context, variant *model.Variant) error
	UpdateVariant(ctx context.Context, variant *model.Variant) error
	DeleteVariant(ctx context.Context, variantID uint, productID uint) error
	FindVariantByID(ctx context.Context, variantID uint, productID uint) (*model.Variant, error)
	CheckVariantExistsBySKU(ctx context.Context, sku string) (bool, error)
	FindVariants(ctx context.Context, where map[string]any, preloads []string) ([]model.Variant, error)
}

type ProductHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Filter(w http.ResponseWriter, r *http.Request)

	// Variants
	CreateVariant(w http.ResponseWriter, r *http.Request)
	UpdateVariant(w http.ResponseWriter, r *http.Request)
	DeleteVariant(w http.ResponseWriter, r *http.Request)
	GetVariant(w http.ResponseWriter, r *http.Request)
}

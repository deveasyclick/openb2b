package product

import (
	"context"
	"errors"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.ProductRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *repository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *repository) Delete(ctx context.Context, ID uint) error {
	res := r.db.WithContext(ctx).Delete(&model.Product{}, ID)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, ID uint) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).First(&product, ID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *repository) Filter(ctx context.Context, opts pagination.Options) ([]model.Product, int64, error) {
	return pagination.Paginate[model.Product](r.db, opts)
}

func (r *repository) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Product, error) {
	var result model.Product

	query := r.db.WithContext(ctx).Model(model.Product{}).Select(fields)

	if where != nil {
		query = query.Where(where)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Variants

func (r *repository) CreateVariant(ctx context.Context, variant *model.Variant) error {
	return r.db.WithContext(ctx).Create(variant).Error
}

func (r *repository) UpdateVariant(ctx context.Context, variant *model.Variant) error {
	return r.db.WithContext(ctx).Save(variant).Error
}

func (r *repository) DeleteVariant(ctx context.Context, variantID uint, productID uint) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND product_id = ?", variantID, productID).
		Delete(&model.Variant{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) FindVariantByID(ctx context.Context, variantID uint, productID uint) (*model.Variant, error) {
	var variant model.Variant
	if err := r.db.WithContext(ctx).
		Where("id = ? AND product_id = ?", variantID, productID).
		First(&variant).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *repository) CheckVariantExistsBySKU(ctx context.Context, sku string) (bool, error) {
	var variant model.Variant
	err := r.db.WithContext(ctx).
		Where("sku = ?", sku).
		First(&variant).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

// WithTx returns a new repository with the given transaction
func (r *repository) WithTx(tx *gorm.DB) interfaces.ProductRepository {
	return &repository{db: tx}
}

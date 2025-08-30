package product

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type service struct {
	repo interfaces.ProductRepository
}

func NewService(repo interfaces.ProductRepository) interfaces.ProductService {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, product *model.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *service) Update(ctx context.Context, product *model.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *service) FindByID(ctx context.Context, ID uint) (*model.Product, error) {
	return s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, nil)
}

func (s *service) Filter(ctx context.Context, opts pagination.Options) ([]model.Product, int64, error) {
	return s.repo.Filter(ctx, opts)
}

func (s *service) Delete(ctx context.Context, ID uint) error {
	return s.repo.Delete(ctx, ID)
}

func (s *service) Exists(ctx context.Context, where map[string]any) (bool, error) {
	p, err := s.repo.FindOneWithFields(ctx, []string{"id"}, where, nil)

	if err != nil {
		return false, err
	}

	return p.ID != 0, nil
}

func (s *service) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Product, error) {
	return s.repo.FindOneWithFields(ctx, fields, where, preloads)
}

// variants

func (s *service) CreateVariant(ctx context.Context, variant *model.Variant) error {
	return s.repo.CreateVariant(ctx, variant)
}

func (s *service) FindVariantByID(ctx context.Context, productID, variantID uint) (*model.Variant, error) {
	return s.repo.FindVariantByID(ctx, variantID, productID)
}

func (s *service) DeleteVariant(ctx context.Context, productID, variantID uint) error {
	return s.repo.DeleteVariant(ctx, variantID, productID)
}

func (s *service) UpdateVariant(ctx context.Context, variant *model.Variant) error {
	return s.repo.UpdateVariant(ctx, variant)
}

func (s *service) CheckVariantExists(ctx context.Context, sku string) (bool, error) {
	return s.repo.CheckVariantExistsBySKU(ctx, sku)
}

func (s *service) WithTx(tx *gorm.DB) interfaces.ProductService {
	return &service{repo: s.repo.WithTx(tx)}
}

package order

import (
	"context"
	"fmt"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type service struct {
	repo           interfaces.OrderRepository
	productService interfaces.ProductService
}

// NewUserService creates a service for orders
func NewService(repo interfaces.OrderRepository, productService interfaces.ProductService) interfaces.OrderService {
	return &service{
		repo:           repo,
		productService: productService,
	}
}

func (s *service) Create(ctx context.Context, DTO dto.CreateOrderDTO, orgId uint) (*model.Order, error) {
	if len(DTO.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}
	// Convert items to non pointer
	item := make([]*dto.CreateOrderItemDTO, len(DTO.Items))
	for i := range DTO.Items {
		item[i] = &DTO.Items[i]
	}

	variantMap, err := s.getVariantMap(ctx, item)
	if err != nil {
		return nil, err
	}

	// Convert DTO to model
	order := DTO.ToModel(variantMap, orgId)

	// Persist order
	if err := s.repo.Create(ctx, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *service) Update(ctx context.Context, order *model.Order, DTO dto.UpdateOrderDTO) error {
	if len(DTO.Items) > 0 {
		variantMap, err := s.getVariantMap(ctx, DTO.Items)
		if err != nil {
			return err
		}
		DTO.ApplyModel(order, &variantMap)
	} else {
		DTO.ApplyModel(order, nil)
	}
	err := s.repo.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FindByID(ctx context.Context, ID uint) (*model.Order, error) {
	return s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, nil)
}

func (s *service) Filter(ctx context.Context, opts pagination.Options) ([]model.Order, int64, error) {
	return s.repo.Filter(ctx, opts)
}

func (s *service) Delete(ctx context.Context, ID uint) error {
	return s.repo.Delete(ctx, ID)
}

func (s *service) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Order, error) {
	return s.repo.FindOneWithFields(ctx, fields, where, preloads)
}

func (s *service) WithTx(tx *gorm.DB) interfaces.OrderService {
	return &service{repo: s.repo.WithTx(tx)}
}

func (s *service) Exists(ctx context.Context, where map[string]any) (bool, error) {
	p, err := s.repo.FindOneWithFields(ctx, []string{"id"}, where, nil)

	if err != nil {
		return false, err
	}

	return p.ID != 0, nil
}

func (s *service) getVariantMap(ctx context.Context, items []*dto.CreateOrderItemDTO) (map[uint]model.Variant, error) {

	var variantIds []uint
	for _, item := range items {
		variantIds = append(variantIds, item.VariantID)
	}

	// Fetch variants
	variants, err := s.productService.FindVariants(ctx, map[string]any{"id": variantIds}, nil)
	if err != nil {
		return nil, err
	}

	// Build map
	variantMap := make(map[uint]model.Variant)
	for _, v := range variants {
		variantMap[v.ID] = v
	}

	// Check for missing variants
	for _, item := range items {
		if _, ok := variantMap[item.VariantID]; !ok {
			return nil, fmt.Errorf("variant %d not found", item.VariantID)
		}
	}

	return variantMap, nil
}

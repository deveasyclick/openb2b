package customer

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type service struct {
	repo interfaces.CustomerRepository
}

func NewService(repo interfaces.CustomerRepository) interfaces.CustomerService {
	return &service{
		repo: repo,
	}
}

func (s *service) Filter(ctx context.Context, opts pagination.Options) ([]model.Customer, int64, error) {
	return s.repo.Filter(ctx, opts)
}

func (s *service) Create(ctx context.Context, customer *model.Customer) error {
	return s.repo.Create(ctx, customer)
}

func (s *service) Update(ctx context.Context, customerId uint, dto *dto.UpdateCustomerDTO) (*model.Customer, error) {
	existingCustomer, err := s.repo.FindByID(ctx, customerId)
	if err != nil {
		return nil, err
	}
	dto.ApplyModel(existingCustomer)
	err = s.repo.Update(ctx, existingCustomer)
	if err != nil {
		return nil, err
	}
	return existingCustomer, nil
}

func (s *service) Delete(ctx context.Context, ID uint) error {
	return s.repo.Delete(ctx, ID)
}

func (s *service) FindByID(ctx context.Context, ID uint, preloads []string) (*model.Customer, error) {
	customer, err := s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, preloads)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *service) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Customer, error) {
	return s.repo.FindOneWithFields(ctx, fields, where, preloads)
}

func (s *service) WithTx(tx *gorm.DB) interfaces.CustomerService {
	return &service{repo: s.repo.WithTx(tx)}
}

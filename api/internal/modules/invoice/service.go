package invoice

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type service struct {
	repo interfaces.InvoiceRepository
	os   interfaces.OrderService
}

func NewService(repo interfaces.InvoiceRepository, os interfaces.OrderService) interfaces.InvoiceService {
	return &service{
		repo: repo,
		os:   os,
	}
}

func (s *service) Filter(ctx context.Context, opts pagination.Options) ([]model.Invoice, int64, error) {
	return s.repo.Filter(ctx, opts)
}

func (s *service) Create(ctx context.Context, orgID uint, dto *dto.CreateInvoiceDTO) (*model.Invoice, error) {
	order, err := s.os.FindByID(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}

	invoice := dto.ToModel(orgID, order)
	err = s.repo.Create(ctx, invoice)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (s *service) Update(ctx context.Context, invoiceId uint, dto *dto.UpdateInvoiceDTO) (*model.Invoice, error) {
	existingInvoice, err := s.repo.FindByID(ctx, invoiceId)
	if err != nil {
		return nil, err
	}
	dto.ApplyModel(existingInvoice)
	err = s.repo.Update(ctx, existingInvoice)
	if err != nil {
		return nil, err
	}
	return existingInvoice, nil
}

func (s *service) Delete(ctx context.Context, ID uint) error {
	return s.repo.Delete(ctx, ID)
}

func (s *service) FindByID(ctx context.Context, ID uint, preloads []string) (*model.Invoice, error) {
	invoice, err := s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, preloads)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (s *service) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Invoice, error) {
	return s.repo.FindOneWithFields(ctx, fields, where, preloads)
}

func (s *service) WithTx(tx *gorm.DB) interfaces.InvoiceService {
	return &service{repo: s.repo.WithTx(tx)}
}

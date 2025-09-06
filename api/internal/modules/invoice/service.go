package invoice

import (
	"context"
	"errors"
	"fmt"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/internal/utils/pdfutil"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type service struct {
	repo   interfaces.InvoiceRepository
	os     interfaces.OrderService
	appCtx *deps.AppContext
}

func NewService(repo interfaces.InvoiceRepository, os interfaces.OrderService, appCtx *deps.AppContext) interfaces.InvoiceService {
	return &service{
		repo:   repo,
		os:     os,
		appCtx: appCtx,
	}
}

func (s *service) Filter(ctx context.Context, opts pagination.Options) ([]model.Invoice, int64, error) {
	return s.repo.Filter(ctx, opts)
}

func (s *service) Create(ctx context.Context, orgID uint, dto *dto.CreateInvoiceDTO) (*model.Invoice, error) {
	order, err := s.os.FindOneWithFields(ctx, nil, map[string]any{"id": dto.OrderID}, []string{"Items", "Customer"})
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

func (s *service) Issue(ctx context.Context, id uint) error {
	invoice, err := s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": id}, []string{"Items", "Order"})
	if err != nil {
		return err
	}

	// If not draft, it's an invalid state for issuing
	if invoice.Status != model.InvoiceStatusDraft {
		return fmt.Errorf("%w: invoice %d is %s", errors.New(apperrors.ErrInvalidInvoiceStatus), id, invoice.Status)
	}

	invoice.Status = model.InvoiceStatusIssued
	err = s.repo.Update(ctx, invoice)
	if err != nil {
		return err
	}

	//TODO: Move email sending to queue
	invCopy := *invoice
	go s.sendInvoiceEmail(&invCopy, s.appCtx.Logger)

	return nil
}

func (s *service) sendInvoiceEmail(invoice *model.Invoice, logger interfaces.Logger) {
	defer func() {
		if r := recover(); r != nil {
			s.appCtx.Logger.Error("panic in sendInvoiceEmail", "err", r)
		}
	}()

	proForma := invoice.Status == model.InvoiceStatusProForma
	pdfBytes, err := pdfutil.GenerateInvoicePDF(invoice, proForma)
	if err != nil {
		logger.Error("failed to generate invoice PDF", "err", err)
		return
	}

	subject := "Your Invoice"
	if proForma {
		subject = "Pro Forma Invoice for Review"
	}

	if err := s.appCtx.Mailer.SendWithAttachment(invoice.CustomerEmail, subject, "Please find attached.", "invoice.pdf", pdfBytes); err != nil {
		logger.Error("failed to send invoice email", "err", err)
		return
	}

	s.appCtx.Logger.Info("invoice email sent", "email", invoice.CustomerEmail)
}

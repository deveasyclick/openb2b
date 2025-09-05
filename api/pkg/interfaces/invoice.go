package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"gorm.io/gorm"
)

type InvoiceHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Filter(w http.ResponseWriter, r *http.Request)
}

type InvoiceService interface {
	Create(ctx context.Context, orgID uint, dto *dto.CreateInvoiceDTO) (*model.Invoice, error)
	Update(ctx context.Context, invoiceId uint, dto *dto.UpdateInvoiceDTO) (*model.Invoice, error)
	Delete(ctx context.Context, ID uint) error
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Invoice, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Invoice, int64, error)
	FindByID(ctx context.Context, ID uint, preloads []string) (*model.Invoice, error)
	WithTx(tx *gorm.DB) InvoiceService
}

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *model.Invoice) error
	Update(ctx context.Context, invoice *model.Invoice) error
	Delete(ctx context.Context, ID uint) error
	FindByID(ctx context.Context, ID uint) (*model.Invoice, error)
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Invoice, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Invoice, int64, error)
	WithTx(tx *gorm.DB) InvoiceRepository
}

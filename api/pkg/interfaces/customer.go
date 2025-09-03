package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"gorm.io/gorm"
)

type CustomerHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Filter(w http.ResponseWriter, r *http.Request)
}

type CustomerService interface {
	Create(ctx context.Context, customer *model.Customer) error
	Update(ctx context.Context, customerId uint, dto *dto.UpdateCustomerDTO) (*model.Customer, error)
	Delete(ctx context.Context, ID uint) error
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Customer, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Customer, int64, error)
	FindByID(ctx context.Context, ID uint, preloads []string) (*model.Customer, error)
	WithTx(tx *gorm.DB) CustomerService
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *model.Customer) error
	Update(ctx context.Context, customer *model.Customer) error
	Delete(ctx context.Context, ID uint) error
	FindByID(ctx context.Context, ID uint) (*model.Customer, error)
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Customer, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Customer, int64, error)
	WithTx(tx *gorm.DB) CustomerRepository
}

package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"gorm.io/gorm"
)

type OrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Filter(w http.ResponseWriter, r *http.Request)
}

type OrderService interface {
	Create(ctx context.Context, DTO dto.CreateOrderDTO, orgId uint) (*model.Order, error)
	Update(ctx context.Context, order *model.Order, dtos dto.UpdateOrderDTO) error
	Delete(ctx context.Context, ID uint) error
	FindByID(ctx context.Context, ID uint) (*model.Order, error)
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Order, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Order, int64, error)
	WithTx(tx *gorm.DB) OrderService
	Exists(ctx context.Context, where map[string]any) (bool, error)
}

type OrderRepository interface {
	Create(ctx context.Context, model *model.Order) error
	Update(ctx context.Context, model *model.Order) error
	FindByID(ctx context.Context, ID uint) (*model.Order, error)
	Filter(ctx context.Context, opts pagination.Options) ([]model.Order, int64, error)
	Delete(ctx context.Context, ID uint) error
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Order, error)
	WithTx(tx *gorm.DB) OrderRepository
}

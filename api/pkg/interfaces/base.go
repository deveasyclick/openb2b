package interfaces

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, model *T) error
	Update(ctx context.Context, model *T) error
	FindByID(ctx context.Context, ID uint) (*T, error)
	Filter(ctx context.Context, opts pagination.Options) ([]T, int64, error)
	Delete(ctx context.Context, ID uint) error
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*T, error)
	WithTx(tx *gorm.DB) Repository[T]
}

type BaseService[T any] interface {
	Create(ctx context.Context, model *T) error
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, ID uint) error
	FindByID(ctx context.Context, ID uint) (*T, error)
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*T, error)
	Filter(ctx context.Context, opts pagination.Options) ([]T, int64, error)
	WithTx(tx *gorm.DB) BaseService[T]
}

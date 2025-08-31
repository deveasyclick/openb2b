package factory

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type BaseService[T any] struct {
	Repo interfaces.Repository[T]
}

func NewService[T any](repo interfaces.Repository[T]) *BaseService[T] {
	return &BaseService[T]{Repo: repo}
}

func (s *BaseService[T]) Create(ctx context.Context, model *T) error {
	return s.Repo.Create(ctx, model)
}

func (s *BaseService[T]) Update(ctx context.Context, model *T) error {
	return s.Repo.Update(ctx, model)
}

func (s *BaseService[T]) FindByID(ctx context.Context, ID uint) (*T, error) {
	return s.Repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, nil)
}

func (s *BaseService[T]) Filter(ctx context.Context, opts pagination.Options) ([]T, int64, error) {
	return s.Repo.Filter(ctx, opts)
}

func (s *BaseService[T]) Delete(ctx context.Context, ID uint) error {
	return s.Repo.Delete(ctx, ID)
}

func (s *BaseService[T]) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*T, error) {
	return s.Repo.FindOneWithFields(ctx, fields, where, preloads)
}

func (s *BaseService[T]) WithTx(tx *gorm.DB) interfaces.BaseService[T] {
	return &BaseService[T]{Repo: s.Repo.WithTx(tx)}
}

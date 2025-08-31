package db

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Filter(ctx context.Context, opts pagination.Options) ([]T, int64, error) {
	return pagination.Paginate[T](r.DB, opts)
}

func (r *BaseRepository[T]) Create(ctx context.Context, model *T) error {
	return r.DB.WithContext(ctx).Create(model).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, model *T) error {
	return r.DB.WithContext(ctx).Save(model).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, ID uint) error {
	var m T
	res := r.DB.WithContext(ctx).Delete(&m, ID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, ID uint) (*T, error) {
	var m T
	err := r.DB.WithContext(ctx).First(&m, ID).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BaseRepository[T]) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*T, error) {
	var result, model T

	query := r.DB.WithContext(ctx).Model(&model).Select(fields)

	if where != nil {
		query = query.Where(where)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *BaseRepository[T]) WithTx(tx *gorm.DB) interfaces.Repository[T] {
	return &BaseRepository[T]{DB: tx}
}

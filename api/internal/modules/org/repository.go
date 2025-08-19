package org

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type orgRepository struct {
	db *gorm.DB
}

func (r *orgRepository) Create(ctx context.Context, org *model.Org) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *orgRepository) Update(ctx context.Context, org *model.Org) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *orgRepository) Delete(ctx context.Context, ID uint) error {
	return r.db.WithContext(ctx).Delete(&model.Org{}, ID).Error
}

func (r *orgRepository) FindByID(ctx context.Context, ID uint) (*model.Org, error) {
	var org model.Org
	err := r.db.WithContext(ctx).First(&org, ID).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *orgRepository) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Org, error) {
	var result model.Org

	query := r.db.WithContext(ctx).Model(model.Org{}).Select(fields)

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

func NewOrgRepository(db *gorm.DB) interfaces.OrgRepository {
	return &orgRepository{
		db: db,
	}
}

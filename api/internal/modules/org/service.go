package org

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type service struct {
	repo interfaces.OrgRepository
}

func NewService(repo interfaces.OrgRepository) interfaces.OrgService {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, org *model.Org) error {
	return s.repo.Create(ctx, org)
}

func (s *service) Update(ctx context.Context, org *model.Org) error {
	return s.repo.Update(ctx, org)
}

func (s *service) Delete(ctx context.Context, ID uint) error {
	return s.repo.Delete(ctx, ID)
}

func (s *service) FindOrg(ctx context.Context, ID uint) (*model.Org, error) {
	return s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, nil)
}

func (s *service) Exists(ctx context.Context, where map[string]any) (bool, error) {
	org, err := s.repo.FindOneWithFields(ctx, []string{"id"}, where, nil)
	if err != nil {
		return false, err
	}

	if org != nil {
		return true, nil
	}

	return false, nil
}

func (s *service) WithTx(tx *gorm.DB) interfaces.OrgService {
	return &service{repo: s.repo.WithTx(tx)}
}

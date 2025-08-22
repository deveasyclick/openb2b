package org

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/apperrors"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type orgService struct {
	repo interfaces.OrgRepository
}

func NewOrgService(repo interfaces.OrgRepository) interfaces.OrgService {
	return &orgService{
		repo: repo,
	}
}

func (s *orgService) Create(ctx context.Context, org *model.Org, userID uint) *apperrors.APIError {
	if err := s.repo.Create(ctx, org); err != nil {
		return &apperrors.APIError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s: id %d", apperrors.ErrCreateOrg, org.ID),
		}
	}

	return nil
}

func (s *orgService) Update(ctx context.Context, org *model.Org) *apperrors.APIError {
	existing, err := s.repo.FindByID(ctx, org.ID)
	if err != nil || existing == nil {
		return &apperrors.APIError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("%s: id %d", apperrors.ErrOrgNotFound, org.ID),
		}
	}

	err = s.repo.Update(ctx, org)
	if err != nil {
		return &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrUpdateOrg,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrUpdateOrg, err),
		}
	}

	return nil
}

func (s *orgService) Delete(ctx context.Context, ID uint) *apperrors.APIError {
	err := s.repo.Delete(ctx, ID)
	if err != nil {
		return &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrDeleteOrg,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrDeleteOrg, err),
		}
	}
	return nil
}

func (s *orgService) FindOrg(ctx context.Context, ID uint) (*model.Org, *apperrors.APIError) {
	org, err := s.repo.FindOneWithFields(ctx, []string{"id"}, map[string]any{"id": ID}, nil)
	if err != nil {
		return nil, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrFindOrg,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrFindOrg, err),
		}
	}
	return org, nil
}

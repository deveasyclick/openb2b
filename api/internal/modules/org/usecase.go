package org

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type CreateOrgUseCase interface {
	Execute(cxt context.Context, input CreateOrgInput) *apperrors.APIError
}

type CreateOrgInput struct {
	Org    *model.Org
	UserID uint
}

func NewCreateOrgUseCase(
	os interfaces.OrgService,
	us interfaces.UserService,
) CreateOrgUseCase {
	return &createOrgUseCase{
		orgService:  os,
		userService: us,
	}
}

type createOrgUseCase struct {
	orgService  interfaces.OrgService
	userService interfaces.UserService
}

func (uc *createOrgUseCase) Execute(ctx context.Context, input CreateOrgInput) *apperrors.APIError {
	err := uc.orgService.Create(ctx, input.Org, input.UserID)
	if err != nil {
		return &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateOrg,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrCreateOrg, err),
		}
	}

	if err := uc.userService.AssignOrg(ctx, input.UserID, input.Org.ID); err != nil {
		return &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrCreateOrg,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrCreateOrg, err),
		}
	}

	return nil
}

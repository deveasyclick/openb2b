package org

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type createOrgUseCase struct {
	orgService  interfaces.OrgService
	userService interfaces.UserService
}

func NewCreateUseCase(
	os interfaces.OrgService,
	us interfaces.UserService,
) interfaces.CreateOrgUseCase {
	return &createOrgUseCase{
		orgService:  os,
		userService: us,
	}
}

func (uc *createOrgUseCase) Execute(ctx context.Context, input types.CreateOrgInput) *apperrors.APIError {
	err := uc.orgService.Create(ctx, input.Org)
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

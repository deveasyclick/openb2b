package org

import (
	"context"
	"errors"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type createOrgUseCase struct {
	orgService   interfaces.OrgService
	userService  interfaces.UserService
	clerkService interfaces.ClerkService
	appCtx       *deps.AppContext
}

func NewCreateUseCase(
	os interfaces.OrgService,
	us interfaces.UserService,
	cs interfaces.ClerkService,
	appCtx *deps.AppContext,
) interfaces.CreateOrgUseCase {
	return &createOrgUseCase{
		orgService:   os,
		userService:  us,
		clerkService: cs,
		appCtx:       appCtx,
	}
}

func (uc *createOrgUseCase) Execute(ctx context.Context, input types.CreateOrgInput) *apperrors.APIError {
	txErr := uc.appCtx.DB.Transaction(func(tx *gorm.DB) error {

		// Wrap services with transactional repos
		orgServiceTx := uc.orgService.WithTx(tx)
		userServiceTx := uc.userService.WithTx(tx)

		//  Create organization
		if err := orgServiceTx.Create(ctx, input.Org); err != nil {
			return err // rollback
		}

		// Assign user to organization
		if err := userServiceTx.AssignOrg(ctx, input.User.ID, input.Org.ID); err != nil {
			return err // rollback
		}

		return nil // commit
	})

	// TODO: return plain error in services
	if txErr != nil {
		var apiErr *apperrors.APIError
		if errors.As(txErr, &apiErr) {
			return apiErr // return your own type
		}

		return &apperrors.APIError{
			Code:        500,
			Message:     "internal server error",
			InternalMsg: txErr.Error(),
		}
	}

	orgID := strconv.FormatUint(uint64(input.Org.ID), 10)
	err := uc.clerkService.SetOrg(ctx, input.User.ClerkID, orgID)
	if err != nil {
		uc.appCtx.Logger.Error("failed to set custom claim org in clerk", "error", err)
		// TODO: Send slerk alerk if this fail or delete org when this fail
		// TODO: Do this later in  a background worker
		return nil
	}

	return nil
}

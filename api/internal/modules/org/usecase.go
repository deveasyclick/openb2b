package org

import (
	"context"
	"errors"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/clerk"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type createOrgUseCase struct {
	orgService   interfaces.OrgService
	userService  interfaces.UserService
	clerkService clerk.Service
	appCtx       *deps.AppContext
}

func NewCreateUseCase(
	os interfaces.OrgService,
	us interfaces.UserService,
	cs clerk.Service,
	appCtx *deps.AppContext,
) interfaces.CreateOrgUseCase {
	return &createOrgUseCase{
		orgService:   os,
		userService:  us,
		clerkService: cs,
		appCtx:       appCtx,
	}
}

func (uc *createOrgUseCase) Execute(ctx context.Context, input types.CreateOrgInput) error {
	txErr := uc.appCtx.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

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

	if txErr != nil {
		return txErr
	}

	orgID := strconv.FormatUint(uint64(input.Org.ID), 10)
	err := uc.clerkService.SetOrg(ctx, input.User.ClerkID, orgID)
	if err != nil {
		// TODO: Send slerk alerk if this fail or delete org when this fail
		// TODO: Do this later in  a background worker
		// Delete org and user org assignment if clerk update failed
		go func() {
			if rec := recover(); rec != nil {
				uc.appCtx.Logger.Error("error recovered from panic", "error", err)
			}
			err := uc.orgService.Delete(ctx, input.Org.ID)
			if err != nil {
				uc.appCtx.Logger.Error("error deleting org", "error", err)
			}
			err = uc.userService.AssignOrg(ctx, input.User.ID, 0)
			if err != nil {
				uc.appCtx.Logger.Error("error deleting org", "error", err)
			}
		}()

		return errors.New("failed to create workspace, setting custom claim in clerk failed")
	}

	return nil
}

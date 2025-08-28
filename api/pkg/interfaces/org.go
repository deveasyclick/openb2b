package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"gorm.io/gorm"
)

type OrgRepository interface {
	Create(ctx context.Context, org *model.Org) error
	Update(ctx context.Context, org *model.Org) error
	Delete(ctx context.Context, ID uint) error
	FindByID(ctx context.Context, ID uint) (*model.Org, error)
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Org, error)
	WithTx(tx *gorm.DB) OrgRepository
}

type OrgService interface {
	Create(ctx context.Context, org *model.Org) *apperrors.APIError
	Update(ctx context.Context, org *model.Org) *apperrors.APIError
	Delete(ctx context.Context, ID uint) *apperrors.APIError
	FindOrg(ctx context.Context, ID uint) (*model.Org, *apperrors.APIError)
	WithTx(tx *gorm.DB) OrgService
	Exists(ctx context.Context, where map[string]any) (bool, *apperrors.APIError)
}

type OrgHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type CreateOrgUseCase interface {
	Execute(cxt context.Context, input types.CreateOrgInput) *apperrors.APIError
}

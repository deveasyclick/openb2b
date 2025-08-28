package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, ID uint) error
	FindByID(ctx context.Context, ID uint) (*model.User, error)
	FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.User, error)
	WithTx(tx *gorm.DB) UserRepository
}

type UserService interface {
	Create(ctx context.Context, user *model.User) *apperrors.APIError
	Update(ctx context.Context, distributor *model.User) *apperrors.APIError
	Delete(ctx context.Context, ID uint) *apperrors.APIError
	FindByEmail(ctx context.Context, email string) (*model.User, *apperrors.APIError)
	FindByID(ctx context.Context, ID uint, preloads []string) (*model.User, *apperrors.APIError)
	AssignOrg(ctx context.Context, userID uint, orgID uint) *apperrors.APIError
	WithTx(tx *gorm.DB) UserService
}

type UserHandler interface {
	GetMe(w http.ResponseWriter, r *http.Request)
}

package interfaces

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) error
	UpdateDistributor(ctx context.Context, ID string, distributor *model.User) error
	DeleteDistributor(ctx context.Context, ID string) error
	GetDistributorByEmail(ctx context.Context, email string) (*model.User, error)
	GetDistributorByID(ctx context.Context, ID string, preloads []string) (*model.User, *apperrors.APIError)
	AssignOrg(ctx context.Context, userID uint, orgID uint) error
}

package clerk

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
)

type Service interface {
	SetOrg(ctx context.Context, clerkUserID string, orgID string) error
	SetRoleAndExternalID(ctx context.Context, userClerkID string, externalID string, role model.Role) error
	DeleteUser(ctx context.Context, userClerkID string) error
}

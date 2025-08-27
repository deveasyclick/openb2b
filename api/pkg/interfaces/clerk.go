package interfaces

import (
	"context"
)

type ClerkService interface {
	SetOrg(ctx context.Context, clerkUserID string, orgID uint) error
	SetExternalID(ctx context.Context, userClerkID string, externalID string) error
	DeleteUser(ctx context.Context, userClerkID string) error
}

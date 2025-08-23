package interfaces

import (
	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
)

type ClerkService interface {
	SetOrg(clerkUserID string, workspaceID string) *apperrors.APIError
	SetExternalID(user *model.User) *apperrors.APIError
}

package clerk

import (
	"context"
	"encoding/json"

	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type service struct {
}

func NewService() interfaces.ClerkService {
	return &service{}
}

func (s *service) SetOrg(ctx context.Context, clerkID string, orgID string) error {
	dataBytes, _ := json.Marshal(map[string]string{
		"org_id": orgID,
	})
	raw := json.RawMessage(dataBytes)
	_, err := user.Update(ctx, clerkID, &user.UpdateParams{PublicMetadata: &raw})

	return err
}

func (s *service) SetExternalID(ctx context.Context, userClerkID string, externalID string) error {
	_, err := user.Update(ctx, userClerkID, &user.UpdateParams{ExternalID: &externalID})

	return err
}

func (s *service) DeleteUser(ctx context.Context, userClerkID string) error {
	_, err := user.Delete(ctx, userClerkID)
	return err
}

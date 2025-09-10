package clerk

import (
	"context"
	"encoding/json"

	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/deveasyclick/openb2b/internal/model"
)

type client struct{}

func New() Service {
	return &client{}
}

func (s *client) SetOrg(ctx context.Context, clerkID string, orgID string) error {
	dataBytes, _ := json.Marshal(map[string]string{
		"org_id": orgID,
	})
	raw := json.RawMessage(dataBytes)
	_, err := user.UpdateMetadata(ctx, clerkID, &user.UpdateMetadataParams{PublicMetadata: &raw})

	return err
}

func (s *client) SetRoleAndExternalID(ctx context.Context, userClerkID string, externalID string, role model.Role) error {
	dataBytes, _ := json.Marshal(map[string]model.Role{
		"role": role,
	})
	raw := json.RawMessage(dataBytes)
	_, err := user.Update(ctx, userClerkID, &user.UpdateParams{ExternalID: &externalID, PublicMetadata: &raw})

	return err
}

func (s *client) DeleteUser(ctx context.Context, userClerkID string) error {
	_, err := user.Delete(ctx, userClerkID)
	return err
}

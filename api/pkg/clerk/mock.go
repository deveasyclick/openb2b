package clerk

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
)

type MockClient struct {
	SetOrgFunc               func(ctx context.Context, clerkUserID string, orgID string) error
	SetRoleAndExternalIDFunc func(ctx context.Context, userClerkID string, externalID string, role model.Role) error
	DeleteUserFunc           func(ctx context.Context, userClerkID string) error
}

func NewMock() Service {
	return &MockClient{}
}
func (m *MockClient) SetOrg(ctx context.Context, clerkID, orgID string) error {
	if m.SetOrgFunc != nil {
		return m.SetOrgFunc(ctx, clerkID, orgID)
	}
	return nil
}

func (m *MockClient) SetRoleAndExternalID(ctx context.Context, clerkID, orgID string, role model.Role) error {
	if m.SetOrgFunc != nil {
		return m.SetRoleAndExternalIDFunc(ctx, clerkID, orgID, role)
	}
	return nil
}

func (m *MockClient) DeleteUser(ctx context.Context, clerkID string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, clerkID)
	}
	return nil
}

package webhook

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/stretchr/testify/require"
)

// --- Mock dependencies ---

type mockUserService struct {
	createFn func(ctx context.Context, u *model.User) *apperrors.APIError
}

func (m *mockUserService) Create(ctx context.Context, u *model.User) *apperrors.APIError {
	return m.createFn(ctx, u)
}
func (m *mockUserService) Update(ctx context.Context, u *model.User) *apperrors.APIError { return nil }
func (m *mockUserService) Delete(ctx context.Context, id uint) *apperrors.APIError       { return nil }
func (m *mockUserService) FindByID(ctx context.Context, id uint, preloads []string) (*model.User, *apperrors.APIError) {
	return nil, nil
}
func (m *mockUserService) AssignOrg(ctx context.Context, userID uint, orgID uint) *apperrors.APIError {
	return nil
}
func (m *mockUserService) FindByEmail(ctx context.Context, email string) (*model.User, *apperrors.APIError) {
	return nil, nil
}

type mockClerkService struct {
	setExternalIDFn func(ctx context.Context, clerkID string, externalID string) error
}

func (m *mockClerkService) SetOrg(ctx context.Context, clerkID string, orgID uint) error {
	return nil
}
func (m *mockClerkService) SetExternalID(ctx context.Context, clerkID string, externalID string) error {
	return m.setExternalIDFn(ctx, clerkID, externalID)
}

type mockLogger struct {
	warnCalled bool
	lastMsg    string
}

func (m *mockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.warnCalled = true
	m.lastMsg = msg
}

func (m *mockLogger) Info(string, ...interface{})  {}
func (m *mockLogger) Debug(string, ...interface{}) {}
func (m *mockLogger) Error(string, ...interface{}) {}
func (m *mockLogger) Fatal(string, ...interface{}) {}
func (m *mockLogger) WithValues(keysAndValues ...interface{}) interfaces.Logger {
	return m
}

// --- Tests ---

func TestHandleEvent_UnknownEvent(t *testing.T) {
	us := &mockUserService{}
	cs := &mockClerkService{}
	appCtx := &deps.AppContext{Logger: &mockLogger{}}

	svc := NewService(us, cs, appCtx)

	err := svc.HandleEvent(context.Background(), &types.WebhookEvent{
		Type: "unknown.event",
		Data: map[string]any{},
	})

	require.Nil(t, err)
}

func TestCreateUser_DecodeError(t *testing.T) {
	us := &mockUserService{}
	cs := &mockClerkService{}
	appCtx := &deps.AppContext{Logger: &mockLogger{}}

	svc := NewService(us, cs, appCtx)

	event := &types.WebhookEvent{
		Type: "user.created",
		Data: map[string]any{"email_addresses": "invalid"},
	}

	err := svc.HandleEvent(context.Background(), event)

	require.Error(t, err)
	require.Equal(t, http.StatusBadRequest, err.Code)
	require.Contains(t, err.Message, "failed to decode ClerkUser")
}

func TestCreateUser_NoEmail(t *testing.T) {
	us := &mockUserService{}
	cs := &mockClerkService{}
	appCtx := &deps.AppContext{Logger: &mockLogger{}}

	svc := NewService(us, cs, appCtx)

	event := &types.WebhookEvent{
		Type: "user.created",
		Data: map[string]any{
			"id":              "clerk_123",
			"first_name":      "John",
			"last_name":       "Doe",
			"email_addresses": []map[string]any{},
		},
	}

	err := svc.HandleEvent(context.Background(), event)

	require.Error(t, err)
	require.Equal(t, http.StatusBadRequest, err.Code)
	require.Contains(t, err.Message, "no email address")
}

func TestCreateUser_CreateFails(t *testing.T) {
	us := &mockUserService{
		createFn: func(ctx context.Context, u *model.User) *apperrors.APIError {
			return &apperrors.APIError{Code: http.StatusInternalServerError, Message: "db fail"}
		},
	}
	cs := &mockClerkService{}
	appCtx := &deps.AppContext{Logger: &mockLogger{}}

	svc := NewService(us, cs, appCtx)

	event := &types.WebhookEvent{
		Type: "user.created",
		Data: map[string]any{
			"id":         "clerk_123",
			"first_name": "John",
			"last_name":  "Doe",
			"email_addresses": []map[string]interface{}{
				{"email_address": "john@example.com"},
			},
		},
	}

	err := svc.HandleEvent(context.Background(), event)

	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, err.Code)
	require.Contains(t, err.Message, "db fail")
}

func TestCreateUser_Success(t *testing.T) {
	us := &mockUserService{
		createFn: func(ctx context.Context, u *model.User) *apperrors.APIError {
			u.ID = 42
			return nil
		},
	}
	cs := &mockClerkService{
		setExternalIDFn: func(ctx context.Context, clerkID string, externalID string) error {
			if clerkID != "clerk_123" || externalID != "42" {
				return errors.New("wrong data passed to SetExternalID")
			}
			return nil
		},
	}
	appCtx := &deps.AppContext{Logger: &mockLogger{}}

	svc := NewService(us, cs, appCtx)

	event := &types.WebhookEvent{
		Type: "user.created",
		Data: map[string]any{
			"id": "clerk_123",
			"email_addresses": []map[string]any{
				{"email_address": "test@example.com"},
			},
			"first_name": "Jane",
			"last_name":  "Doe",
		},
	}

	err := svc.HandleEvent(context.Background(), event)

	require.Nil(t, err)
}

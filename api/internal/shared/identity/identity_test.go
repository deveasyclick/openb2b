package identity

import (
	"context"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
)

// --- helpers ---

// fakeClaims returns a context with clerk.SessionClaims injected
func fakeClaimsCtx(subject string, custom any) context.Context {
	claims := &clerk.SessionClaims{
		RegisteredClaims: clerk.RegisteredClaims{
			Subject: subject,
		},
		Custom: custom,
	}
	ctx := clerk.ContextWithSessionClaims(context.Background(), claims)
	return ctx
}

// --- tests ---

func TestGetCustomClaims(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		want      *CustomSessionClaims
		expectErr bool
	}{
		{
			name:      "no claims in context",
			ctx:       context.Background(),
			want:      nil,
			expectErr: true,
		},
		{
			name: "claims present but wrong type",
			ctx:  fakeClaimsCtx("user-1", "not-a-struct"),
			want: nil, expectErr: true,
		},
		{
			name: "valid custom claims",
			ctx: fakeClaimsCtx("user-1", &CustomSessionClaims{
				UserID: "42",
				OrgID:  "7",
			}),
			want: &CustomSessionClaims{UserID: "42", OrgID: "7"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCustomClaims(tt.ctx)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.UserID != tt.want.UserID || got.OrgID != tt.want.OrgID {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestClerkIDFromContext(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		wantID string
		wantOK bool
	}{
		{
			name:   "no claims",
			ctx:    context.Background(),
			wantID: "",
			wantOK: false,
		},
		{
			name:   "with claims",
			ctx:    fakeClaimsCtx("user-123", &CustomSessionClaims{}),
			wantID: "user-123",
			wantOK: false, // function always returns false as second arg
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotOK := ClerkIDFromContext(tt.ctx)
			if gotID != tt.wantID {
				t.Errorf("gotID = %q, want %q", gotID, tt.wantID)
			}
			if gotOK != tt.wantOK {
				t.Errorf("gotOK = %v, want %v", gotOK, tt.wantOK)
			}
		})
	}
}

func TestUserFromContext(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		want      *ContextUser
		expectErr bool
	}{
		{
			name:      "no claims in context",
			ctx:       context.Background(),
			want:      nil,
			expectErr: true,
		},
		{
			name: "invalid userID",
			ctx: fakeClaimsCtx("user-x", &CustomSessionClaims{
				UserID: "abc",
				OrgID:  "1",
			}),
			want:      nil,
			expectErr: true,
		},
		{
			name: "invalid orgID",
			ctx: fakeClaimsCtx("user-x", &CustomSessionClaims{
				UserID: "10",
				OrgID:  "abc",
			}),
			want:      nil,
			expectErr: true,
		},
		{
			name: "valid IDs",
			ctx: fakeClaimsCtx("user-x", &CustomSessionClaims{
				UserID: "10",
				OrgID:  "20",
			}),
			want: &ContextUser{ID: 10, Org: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserFromContext(tt.ctx)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.want.ID || got.Org != tt.want.Org {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

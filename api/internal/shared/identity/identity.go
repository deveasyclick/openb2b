package identity

import (
	"context"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	parseuint "github.com/deveasyclick/openb2b/internal/utils/parseUint"
)

type CustomSessionClaims struct {
	OrgID   string `json:"org_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	ClerkId string
}

type ContextUser struct {
	ID      uint
	Org     uint
	ClerkID string
}

func GetCustomClaims(ctx context.Context) (*CustomSessionClaims, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no session claims found in context")
	}
	customClaims, ok := claims.Custom.(*CustomSessionClaims)
	if !ok {
		return nil, fmt.Errorf("invalid or missing custom claims")
	}

	return &CustomSessionClaims{
		UserID:  customClaims.UserID,
		OrgID:   customClaims.OrgID,
		ClerkId: claims.Subject,
	}, nil
}

func ClerkIDFromContext(ctx context.Context) (string, bool) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return "", false
	}
	return claims.Subject, false
}

func UserFromContext(ctx context.Context) (*ContextUser, error) {
	claims, err := GetCustomClaims(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := parseuint.ParseUint(claims.UserID, "user ID")
	if err != nil {
		return nil, err
	}

	user := &ContextUser{
		ID:      userID,
		ClerkID: claims.ClerkId,
	}

	return user, nil
}

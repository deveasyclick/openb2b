package types

import (
	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/identity"
)

type CreateOrgInput struct {
	Org  *model.Org
	User *identity.ContextUser
}

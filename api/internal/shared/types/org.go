package types

import "github.com/deveasyclick/openb2b/internal/model"

type CreateOrgInput struct {
	Org    *model.Org
	UserID uint
}

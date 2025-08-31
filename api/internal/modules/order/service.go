package order

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/factory"
	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type service struct {
	*factory.BaseService[model.Order]
}

// NewUserService creates a service for orders
func NewService(repo interfaces.Repository[model.Order]) *service {
	return &service{
		BaseService: factory.NewService(repo),
	}
}

func (s *service) Exists(ctx context.Context, where map[string]any) (bool, error) {
	p, err := s.Repo.FindOneWithFields(ctx, []string{"id"}, where, nil)

	if err != nil {
		return false, err
	}

	return p.ID != 0, nil
}

package user

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type service struct {
	repo interfaces.UserRepository
}

func NewService(repo interfaces.UserRepository) interfaces.UserService {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *service) Update(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, ID uint) error {
	return s.repo.Delete(ctx, ID)
}

func (s *service) FindByID(ctx context.Context, ID uint, preloads []string) (*model.User, error) {
	return s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, preloads)
}

func (s *service) AssignOrg(ctx context.Context, userID uint, orgID uint) error {
	// Associate the workspace with the distributor
	user := &model.User{OrgID: &orgID}
	user.ID = userID
	return s.repo.Update(ctx, user)
}

func (s *service) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.FindOneWithFields(ctx, nil, map[string]any{"email": email}, nil)
}

func (s *service) WithTx(tx *gorm.DB) interfaces.UserService {
	return &service{repo: s.repo.WithTx(tx)}
}

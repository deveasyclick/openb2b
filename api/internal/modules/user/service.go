package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

type service struct {
	repo interfaces.UserRepository
}

func NewService(repo interfaces.UserRepository) interfaces.UserService {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, user *model.User) *apperrors.APIError {
	if err := s.repo.Create(ctx, user); err != nil {
		return &apperrors.APIError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s: id %d", apperrors.ErrCreateUser, user.ID),
		}
	}

	return nil
}

func (s *service) Update(ctx context.Context, user *model.User) *apperrors.APIError {
	existing, err := s.repo.FindByID(ctx, user.ID)
	if err != nil || existing == nil {
		return &apperrors.APIError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("%s: id %d", apperrors.ErrUserNotFound, user.ID),
		}
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrUpdateUser,
			InternalMsg: fmt.Sprintf("%s: error %s", apperrors.ErrUpdateUser, err),
		}
	}

	return nil
}

func (s *service) Delete(ctx context.Context, ID uint) *apperrors.APIError {
	existing, err := s.repo.FindOneWithFields(ctx, []string{"id"}, map[string]any{"id": ID}, nil)
	if err != nil || existing == nil {
		return &apperrors.APIError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("%s: id %d", apperrors.ErrUserNotFound, ID),
		}
	}

	err = s.repo.Delete(ctx, ID)
	if err != nil {
		return &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrDeleteUser,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrDeleteUser, err),
		}
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, ID uint, preloads []string) (*model.User, *apperrors.APIError) {
	user, err := s.repo.FindOneWithFields(ctx, nil, map[string]any{"id": ID}, preloads)
	if err != nil {
		return nil, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrFindUser,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrFindUser, err),
		}
	}
	return user, nil
}

func (s *service) AssignOrg(ctx context.Context, userID uint, orgID uint) *apperrors.APIError {
	// Associate the workspace with the distributor
	user := &model.User{OrgID: orgID}
	user.ID = userID
	err := s.repo.Update(ctx, user)
	if err != nil {
		return &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrFindUser,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrFindUser, err),
		}
	}

	return nil
}

func (s *service) FindByEmail(ctx context.Context, email string) (*model.User, *apperrors.APIError) {
	distributor, err := s.repo.FindOneWithFields(ctx, nil, map[string]any{"email": email}, nil)
	if err != nil {
		return nil, &apperrors.APIError{
			Code:        http.StatusInternalServerError,
			Message:     apperrors.ErrFindUser,
			InternalMsg: fmt.Sprintf("%s: %s", apperrors.ErrFindUser, err),
		}
	}

	return distributor, nil
}

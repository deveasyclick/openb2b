package user

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	createFn            func(ctx context.Context, m *model.User) error
	findByIDFn          func(ctx context.Context, id uint) (*model.User, error)
	updateFn            func(ctx context.Context, m *model.User) error
	deleteFn            func(ctx context.Context, id uint) error
	findOneWithFieldsFn func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error)
}

func (m *mockUserRepo) Create(ctx context.Context, u *model.User) error {
	return m.createFn(ctx, u)
}
func (m *mockUserRepo) FindByID(ctx context.Context, id uint) (*model.User, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockUserRepo) Update(ctx context.Context, u *model.User) error {
	return m.updateFn(ctx, u)
}
func (m *mockUserRepo) Delete(ctx context.Context, id uint) error {
	return m.deleteFn(ctx, id)
}
func (m *mockUserRepo) FindOneWithFields(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
	return m.findOneWithFieldsFn(ctx, fields, cond, preload)
}

func (m *mockUserRepo) WithTx(tx *gorm.DB) interfaces.UserRepository {
	return m
}

func TestService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{
			createFn: func(ctx context.Context, u *model.User) error { return nil },
		}
		svc := NewService(repo)
		user := &model.User{}
		user.ID = 1
		err := svc.Create(context.Background(), user)
		require.Nil(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		repo := &mockUserRepo{
			createFn: func(ctx context.Context, u *model.User) error { return errors.New("db fail") },
		}
		svc := NewService(repo)
		user := &model.User{}
		user.ID = 1
		err := svc.Create(context.Background(), user)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.Message, apperrors.ErrCreateUser)
		require.Contains(t, err.InternalMsg, "db fail")
	})
}

func TestService_Update(t *testing.T) {
	u := &model.User{}
	u.ID = 2

	t.Run("not found", func(t *testing.T) {
		repo := &mockUserRepo{
			findByIDFn: func(ctx context.Context, id uint) (*model.User, error) {
				return nil, errors.New("not found")
			},
		}
		svc := NewService(repo)
		err := svc.Update(context.Background(), u)
		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, err.Code)
		require.Contains(t, err.Message, apperrors.ErrUserNotFound)
	})

	t.Run("update fails", func(t *testing.T) {
		repo := &mockUserRepo{
			findByIDFn: func(ctx context.Context, id uint) (*model.User, error) {
				return u, nil
			},
			updateFn: func(ctx context.Context, u *model.User) error {
				return errors.New("db fail")
			},
		}
		svc := NewService(repo)
		err := svc.Update(context.Background(), u)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.Message, apperrors.ErrUpdateUser)
		require.Contains(t, err.InternalMsg, "db fail")
	})

	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{
			findByIDFn: func(ctx context.Context, id uint) (*model.User, error) {
				return u, nil
			},
			updateFn: func(ctx context.Context, u *model.User) error { return nil },
		}
		svc := NewService(repo)
		err := svc.Update(context.Background(), u)
		require.Nil(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		repo := &mockUserRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
				return nil, errors.New("not found")
			},
		}
		svc := NewService(repo)
		err := svc.Delete(context.Background(), 3)
		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("delete fails", func(t *testing.T) {
		user := &model.User{}
		user.ID = 3
		repo := &mockUserRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
				return user, nil
			},
			deleteFn: func(ctx context.Context, id uint) error {
				return errors.New("delete fail")
			},
		}
		svc := NewService(repo)
		err := svc.Delete(context.Background(), 3)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.InternalMsg, "delete fail")
	})

	t.Run("success", func(t *testing.T) {
		user := &model.User{}
		user.ID = 3
		repo := &mockUserRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
				return user, nil
			},
			deleteFn: func(ctx context.Context, id uint) error { return nil },
		}
		svc := NewService(repo)
		err := svc.Delete(context.Background(), 3)
		require.Nil(t, err)
	})
}

func TestService_FindByID(t *testing.T) {
	t.Run("db error", func(t *testing.T) {
		repo := &mockUserRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
				return nil, errors.New("db fail")
			},
		}
		svc := NewService(repo)
		_, err := svc.FindByID(context.Background(), 4, nil)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.InternalMsg, "db fail")
	})

	t.Run("success", func(t *testing.T) {
		user := &model.User{}
		user.ID = 4
		repo := &mockUserRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
				return user, nil
			},
		}
		svc := NewService(repo)
		got, err := svc.FindByID(context.Background(), 4, nil)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, uint(4), got.ID)
	})
}

func TestService_AssignOrg(t *testing.T) {
	t.Run("update fails", func(t *testing.T) {
		repo := &mockUserRepo{
			updateFn: func(ctx context.Context, u *model.User) error {
				return errors.New("update fail")
			},
		}
		svc := NewService(repo)
		err := svc.AssignOrg(context.Background(), 5, 10)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.InternalMsg, "update fail")
	})

	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{
			updateFn: func(ctx context.Context, u *model.User) error { return nil },
		}
		svc := NewService(repo)
		err := svc.AssignOrg(context.Background(), 5, 10)
		require.Nil(t, err)
	})
}

func TestService_FindByEmail(t *testing.T) {
	t.Run("db error", func(t *testing.T) {
		repo := &mockUserRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
				return nil, errors.New("db fail")
			},
		}
		svc := NewService(repo)
		_, err := svc.FindByEmail(context.Background(), "test@example.com")
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.InternalMsg, "db fail")
	})

	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.User, error) {
				return &model.User{Email: "test@example.com"}, nil
			},
		}
		svc := NewService(repo)
		got, err := svc.FindByEmail(context.Background(), "test@example.com")
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, "test@example.com", got.Email)
	})
}

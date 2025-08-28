package org

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockRepo struct {
	createFn            func(ctx context.Context, m *model.Org) error
	findByIDFn          func(ctx context.Context, id uint) (*model.Org, error)
	updateFn            func(ctx context.Context, m *model.Org) error
	deleteFn            func(ctx context.Context, id uint) error
	findOneWithFieldsFn func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.Org, error)
}

func (m *mockRepo) Create(ctx context.Context, o *model.Org) error {
	return m.createFn(ctx, o)
}
func (m *mockRepo) FindByID(ctx context.Context, id uint) (*model.Org, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockRepo) Update(ctx context.Context, o *model.Org) error {
	return m.updateFn(ctx, o)
}
func (m *mockRepo) Delete(ctx context.Context, id uint) error {
	return m.deleteFn(ctx, id)
}
func (m *mockRepo) FindOneWithFields(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.Org, error) {
	return m.findOneWithFieldsFn(ctx, fields, cond, preload)
}

func (m *mockRepo) WithTx(tx *gorm.DB) interfaces.OrgRepository {
	return m
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name         string
		repo         *mockRepo
		wantError    bool
		wantMsgMatch string
	}{
		{
			name: "create success",
			repo: &mockRepo{
				createFn: func(ctx context.Context, m *model.Org) error { return nil },
			},
			wantError: false,
		},
		{
			name: "create fails",
			repo: &mockRepo{
				createFn: func(ctx context.Context, m *model.Org) error {
					return errors.New("db error")
				},
			},
			wantError:    true,
			wantMsgMatch: "id 1", // matches the formatting in service
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.repo)
			org := &model.Org{}
			org.ID = 1

			err := svc.Create(context.Background(), org)

			if tt.wantError {
				require.NotNil(t, err)
				require.Equal(t, http.StatusInternalServerError, err.Code)
				require.Contains(t, err.Message, tt.wantMsgMatch)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	orgModel := &model.Org{}
	orgModel.ID = 2

	t.Run("not found", func(t *testing.T) {
		repo := &mockRepo{
			findByIDFn: func(ctx context.Context, id uint) (*model.Org, error) {
				return nil, errors.New("not found")
			},
		}
		svc := NewService(repo)

		err := svc.Update(context.Background(), orgModel)

		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.Code)
		require.Contains(t, err.Message, "not found")
	})

	t.Run("update error", func(t *testing.T) {
		repo := &mockRepo{
			findByIDFn: func(ctx context.Context, id uint) (*model.Org, error) {
				return orgModel, nil
			},
			updateFn: func(ctx context.Context, o *model.Org) error {
				return errors.New("db fail")
			},
		}
		svc := NewService(repo)

		err := svc.Update(context.Background(), orgModel)

		require.NotNil(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.InternalMsg, "db fail")
	})

	t.Run("update success", func(t *testing.T) {
		repo := &mockRepo{
			findByIDFn: func(ctx context.Context, id uint) (*model.Org, error) {
				return orgModel, nil
			},
			updateFn: func(ctx context.Context, o *model.Org) error { return nil },
		}
		svc := NewService(repo)

		err := svc.Update(context.Background(), orgModel)

		require.Nil(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("delete error", func(t *testing.T) {
		repo := &mockRepo{
			deleteFn: func(ctx context.Context, id uint) error {
				return errors.New("delete fail")
			},
		}
		svc := NewService(repo)

		err := svc.Delete(context.Background(), 5)

		require.NotNil(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.InternalMsg, "delete fail")
	})

	t.Run("delete success", func(t *testing.T) {
		repo := &mockRepo{
			deleteFn: func(ctx context.Context, id uint) error { return nil },
		}
		svc := NewService(repo)

		err := svc.Delete(context.Background(), 5)

		require.Nil(t, err)
	})
}

func TestService_FindOrg(t *testing.T) {
	orgModel := &model.Org{}
	orgModel.ID = 3

	t.Run("find error", func(t *testing.T) {
		repo := &mockRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.Org, error) {
				return nil, errors.New("db fail")
			},
		}
		svc := NewService(repo)

		_, err := svc.FindOrg(context.Background(), 3)

		require.NotNil(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.InternalMsg, "db fail")
	})

	t.Run("find success", func(t *testing.T) {
		repo := &mockRepo{
			findOneWithFieldsFn: func(ctx context.Context, fields []string, cond map[string]any, preload []string) (*model.Org, error) {
				return orgModel, nil
			},
		}
		svc := NewService(repo)

		got, err := svc.FindOrg(context.Background(), 3)

		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, uint(3), got.ID)
	})
}

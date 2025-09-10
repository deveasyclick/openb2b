package user

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Updates(user).Error
}

func (r *repository) Delete(ctx context.Context, ID uint) error {
	res := r.db.WithContext(ctx).Delete(&model.User{}, ID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, ID uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, ID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.User, error) {
	var result model.User

	query := r.db.WithContext(ctx).Model(model.User{}).Select(fields)

	if where != nil {
		query = query.Where(where)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// WithTx returns a new repository with the given transaction
func (r *repository) WithTx(tx *gorm.DB) interfaces.UserRepository {
	return &repository{db: tx}
}

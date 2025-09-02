package customer

import (
	"context"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/pagination"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.CustomerRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Filter(ctx context.Context, opts pagination.Options) ([]model.Customer, int64, error) {
	return pagination.Paginate[model.Customer](r.db, opts)
}

func (r *repository) Create(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *repository) Update(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Updates(customer).Error
}

func (r *repository) Delete(ctx context.Context, ID uint) error {
	res := r.db.WithContext(ctx).Delete(&model.Customer{}, ID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, ID uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).First(&customer, ID).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *repository) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Customer, error) {
	var result model.Customer

	query := r.db.WithContext(ctx).Model(model.Customer{}).Select(fields)

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
func (r *repository) WithTx(tx *gorm.DB) interfaces.CustomerRepository {
	return &repository{db: tx}
}

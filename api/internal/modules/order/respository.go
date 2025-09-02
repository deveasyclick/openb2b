package order

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

func NewRepository(database *gorm.DB) interfaces.OrderRepository {
	return &repository{
		db: database,
	}
}

func (r *repository) Filter(ctx context.Context, opts pagination.Options) ([]model.Order, int64, error) {
	return pagination.Paginate[model.Order](r.db, opts)
}

func (r *repository) Create(ctx context.Context, model *model.Order) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository) Update(ctx context.Context, model *model.Order) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *repository) Delete(ctx context.Context, ID uint) error {
	res := r.db.WithContext(ctx).Delete(&model.Order{}, ID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, ID uint) (*model.Order, error) {
	var m model.Order
	err := r.db.WithContext(ctx).First(&m, ID).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Order, error) {
	var result model.Order

	query := r.db.WithContext(ctx).Model(&model.Order{}).Select(fields)

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

func (r *repository) WithTx(tx *gorm.DB) interfaces.OrderRepository {
	return &repository{db: tx}
}

// We select all the fields that are needed to update the order so gorm doesn't insert new order due to embedded struct in order
func (r *repository) UpdateAndReplace(cxt context.Context, ID uint, order *model.Order) error {
	return r.db.WithContext(cxt).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Order{}).Where("id = ?", ID).Select(
			"status", "notes", "customer_id",
			"total_price", "total_weight", "discount",
			"tax", "delivery_status", "delivery_transport_fare",
			"discount_amount", "discount_type",
			"delivery_address_zip", "delivery_address_state", "delivery_address_city", "delivery_address_country",
			"delivery_address_address",
		).Updates(order).Error; err != nil {
			return err
		}
		if err := tx.Model(order).Association("Items").Replace(order.Items); err != nil {
			return err
		}
		return nil
	})
}

package invoice

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

func NewRepository(db *gorm.DB) interfaces.InvoiceRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Filter(ctx context.Context, opts pagination.Options) ([]model.Invoice, int64, error) {
	return pagination.Paginate[model.Invoice](r.db, opts)
}

func (r *repository) Create(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).Create(invoice).Error
}

func (r *repository) Update(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).Updates(invoice).Error
}

func (r *repository) Delete(ctx context.Context, ID uint) error {
	res := r.db.WithContext(ctx).Delete(&model.Invoice{}, ID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, ID uint) (*model.Invoice, error) {
	var invoice model.Invoice
	err := r.db.WithContext(ctx).First(&invoice, ID).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *repository) FindOneWithFields(ctx context.Context, fields []string, where map[string]any, preloads []string) (*model.Invoice, error) {
	var result model.Invoice

	query := r.db.WithContext(ctx).Model(model.Invoice{}).Select(fields)

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
func (r *repository) WithTx(tx *gorm.DB) interfaces.InvoiceRepository {
	return &repository{db: tx}
}

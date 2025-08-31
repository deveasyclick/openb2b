package order

import (
	"github.com/deveasyclick/openb2b/internal/db"
	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	*db.BaseRepository[model.Order]
}

func NewRepository(database *gorm.DB) interfaces.OrderRepository {
	return &repository{
		BaseRepository: db.NewBaseRepository[model.Order](database),
	}
}

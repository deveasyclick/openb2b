package order

import (
	"github.com/deveasyclick/openb2b/internal/factory"
	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	*factory.BaseRepository[model.Order]
}

func NewRepository(database *gorm.DB) interfaces.OrderRepository {
	return &repository{
		BaseRepository: factory.NewBaseRepository[model.Order](database),
	}
}

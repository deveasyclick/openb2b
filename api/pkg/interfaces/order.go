package interfaces

import "github.com/deveasyclick/openb2b/internal/model"

type OrderRepository interface {
	Repository[model.Order]
}

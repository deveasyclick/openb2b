package interfaces

import (
	"context"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
)

type OrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Filter(w http.ResponseWriter, r *http.Request)
}

type OrderService interface {
	BaseService[model.Order]
	Exists(ctx context.Context, where map[string]any) (bool, error)
}
type OrderRepository interface {
	Repository[model.Order]
}

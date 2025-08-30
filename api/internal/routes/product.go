package routes

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func registerProductRoutes(router chi.Router, productHandler interfaces.ProductHandler) {

	router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.Filter)

		r.Post("/", productHandler.Create)

		r.Get("/{id}", productHandler.Get)

		r.Patch("/{id}", productHandler.Update)

		r.Delete("/{id}", productHandler.Delete)

		r.Route("/{productId}/variants", func(r chi.Router) {
			r.Post("/", productHandler.CreateVariant)
			r.Patch("/{id}", productHandler.UpdateVariant)
			r.Delete("/{id}", productHandler.DeleteVariant)
			r.Get("/{id}", productHandler.GetVariant)
		})

	})
}

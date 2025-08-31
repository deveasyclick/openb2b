package routes

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func registerOrderRoutes(router chi.Router, orderHandler interfaces.OrderHandler) {

	router.Route("/orders", func(r chi.Router) {
		r.Get("/", orderHandler.Filter)

		r.Post("/", orderHandler.Create)

		r.Get("/{id}", orderHandler.Get)

		r.Patch("/{id}", orderHandler.Update)

		r.Delete("/{id}", orderHandler.Delete)

	})
}

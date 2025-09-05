package routes

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func registerInvoiceRoutes(router chi.Router, handler interfaces.InvoiceHandler) {
	router.Route("/invoices", func(r chi.Router) {
		r.Get("/", handler.Filter)

		r.Post("/", handler.Create)

		r.Get("/{id}", handler.Get)

		r.Patch("/{id}", handler.Update)

		r.Delete("/{id}", handler.Delete)
	})
}

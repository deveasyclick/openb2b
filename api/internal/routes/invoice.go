package routes

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func registerInvoiceRoutes(router chi.Router, handler interfaces.InvoiceHandler) {
	router.Route("/invoices", func(r chi.Router) {
		r.Get("/", handler.Filter)

		r.Post("/", handler.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.Get)
			r.Put("/", handler.Update)
			r.Delete("/", handler.Delete)

			r.Post("/issue", handler.Issue)
		})
	})
}

package routes

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func registerOrgRoutes(router chi.Router, handler interfaces.OrgHandler) {

	router.Route("/orgs", func(r chi.Router) {
		r.Post("/", handler.Create)

		r.Get("/{id}", handler.Get)

		r.Put("/{id}", handler.Update)

		r.Delete("/{id}", handler.Delete)
	})
}

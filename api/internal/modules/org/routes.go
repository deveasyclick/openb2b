package org

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func RegisterRoutes(router chi.Router, orgHandler interfaces.OrgHandler) {

	router.Route("/orgs", func(r chi.Router) {
		r.Post("/", orgHandler.Create)

		r.Get("/{id}", orgHandler.Get)

		r.Patch("/{id}", orgHandler.Update)

		r.Delete("/{id}", orgHandler.Delete)
	})
}

package routes

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func registerUserRoutes(router chi.Router, handler interfaces.UserHandler) {

	router.Route("/users", func(r chi.Router) {
		r.Post("/me", handler.GetMe)
	})
}

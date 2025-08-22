package routes

import (
	"github.com/deveasyclick/openb2b/internal/modules/org"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/go-chi/chi"
)

func registerRoutes(router chi.Router, app *deps.AppContext) {
	orgRepository := org.NewOrgRepository(app.DB)
	orgService := org.NewOrgService(orgRepository)
	// TODO: Implement user service and replace nil with userService
	createOrgUseCase := org.NewCreateOrgUseCase(orgService, nil)
	orgHandler := org.NewOrgHandler(orgService, createOrgUseCase)

	router.Route("/orgs", func(r chi.Router) {
		r.Post("/", orgHandler.Create)

		r.Get("/{id}", orgHandler.Get)

		r.Put("/{id}", orgHandler.Update)

		r.Delete("/{id}", orgHandler.Delete)
	})
}

package routes

import (
	"time"

	"github.com/deveasyclick/openb2b/internal/modules/clerk"
	"github.com/deveasyclick/openb2b/internal/modules/org"
	"github.com/deveasyclick/openb2b/internal/modules/user"
	"github.com/deveasyclick/openb2b/internal/modules/webhook"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func Register(r chi.Router, appCtx *deps.AppContext) {
	r.Use(chiMiddleware.RequestID) // Adds a unique request ID
	r.Use(chiMiddleware.RealIP)    // Gets the real IP from X-Forwarded-For
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	// Enable rate limiter of 100 requests per minute per IP
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(chiMiddleware.Heartbeat("/ping"))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// User
	userRepository := user.NewRepository(appCtx.DB)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService, appCtx)

	// Clerk
	clerkService := clerk.NewService()

	// Webhook
	webhookService := webhook.NewService(userService, clerkService, appCtx)
	webhookHandler := webhook.NewHandler(webhookService, appCtx)

	// Org
	orgRepository := org.NewRepository(appCtx.DB)
	orgService := org.NewService(orgRepository)
	createOrgUseCase := org.NewCreateUseCase(orgService, userService)
	orgHandler := org.NewHandler(orgService, createOrgUseCase, appCtx)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))

		// Public routes
		r.Group(func(r chi.Router) {

		})

		// Private routes
		r.Group(func(r chi.Router) {
			registerOrgRoutes(r, orgHandler)
			registerUserRoutes(r, userHandler)
			registerWebhookRoutes(r, webhookHandler, appCtx)
		})

	})
}

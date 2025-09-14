package routes

import (
	"fmt"
	"net/url"
	"time"

	"github.com/deveasyclick/openb2b/docs"
	"github.com/deveasyclick/openb2b/internal/modules/customer"
	"github.com/deveasyclick/openb2b/internal/modules/invoice"
	"github.com/deveasyclick/openb2b/internal/modules/order"
	"github.com/deveasyclick/openb2b/internal/modules/org"
	"github.com/deveasyclick/openb2b/internal/modules/product"
	"github.com/deveasyclick/openb2b/internal/modules/user"
	"github.com/deveasyclick/openb2b/internal/modules/webhook"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/pkg/clerk"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	swagger "github.com/swaggo/http-swagger"
)

func Register(r chi.Router, appCtx *deps.AppContext, middleware interfaces.Middleware, clerkService clerk.Service) {
	r.Use(chiMiddleware.RequestID) // Adds a unique request ID
	r.Use(chiMiddleware.RealIP)    // Gets the real IP from X-Forwarded-For
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	// Enable rate limiter of 100 requests per minute per IP
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(chiMiddleware.Heartbeat("/ping"))
	r.Use(middleware.Recover(appCtx.Logger))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
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

	// Webhook
	webhookService := webhook.NewService(userService, clerkService, appCtx)
	webhookHandler := webhook.NewHandler(webhookService, appCtx)

	// Org
	orgRepository := org.NewRepository(appCtx.DB)
	orgService := org.NewService(orgRepository)
	createOrgUseCase := org.NewCreateUseCase(orgService, userService, clerkService, appCtx)
	orgHandler := org.NewHandler(orgService, createOrgUseCase, appCtx)

	// Product
	productRepository := product.NewRepository(appCtx.DB)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService, appCtx)

	// Order
	orderRepository := order.NewRepository(appCtx.DB)
	orderService := order.NewService(orderRepository, productService)
	orderHandler := order.NewHandler(orderService, appCtx)

	// Customer
	customerRepository := customer.NewRepository(appCtx.DB)
	customerService := customer.NewService(customerRepository)
	customerHandler := customer.NewHandler(customerService, appCtx)

	// Invoice
	invoiceRepository := invoice.NewRepository(appCtx.DB)
	invoiceService := invoice.NewService(invoiceRepository, orderService, appCtx)
	invoiceHandler := invoice.NewHandler(invoiceService, appCtx)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))

		// Public routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyWebhook())
			registerWebhookRoutes(r, webhookHandler, appCtx)
		})

		// Private routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.ValidateJWT())
			org.RegisterRoutes(r, orgHandler)
			registerUserRoutes(r, userHandler)
			registerProductRoutes(r, productHandler)
			registerOrderRoutes(r, orderHandler)
			registerCustomerRoutes(r, customerHandler)
			registerInvoiceRoutes(r, invoiceHandler)
		})
	})

	if appCtx.Config.Env == "development" {
		parsedURL, err := url.Parse(appCtx.Config.AppURL)
		if err != nil {
			appCtx.Logger.Warn("failed to parse app url", "err", err)
		}

		docs.SwaggerInfo.Host = parsedURL.Host
		r.Get("/swagger/*", swagger.Handler(
			swagger.URL(fmt.Sprintf("%s/swagger/doc.json", appCtx.Config.AppURL)),
		))
	}
}

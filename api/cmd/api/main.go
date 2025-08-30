// @title OpenB2B API
// @version 1.0
// @description Open-source multi-tenant ordering & invoicing platform API.
// @termsOfService http://openb2b.com/terms/

// @contact.name API Support
// @contact.url http://openb2b.com/support
// @contact.email support@openb2b.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/deveasyclick/openb2b/internal/config"
	"github.com/deveasyclick/openb2b/internal/db"
	"github.com/deveasyclick/openb2b/internal/middleware"
	"github.com/deveasyclick/openb2b/internal/routes"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/pkg/logger"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	logger := logger.New(os.Getenv("ENV"))

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal("failed to load config", "err", err)
	}

	clerk.SetKey(cfg.ClerkSecret)

	dbConn := db.New(db.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		Env:      cfg.Env,
	}, logger)

	appCtx := &deps.AppContext{
		DB:     dbConn,
		Config: cfg,
		Logger: logger,
		Cache:  nil,
	}

	middlewares := middleware.New()

	routes.Register(r, appCtx, middlewares)

	port := cfg.Port
	if port == 0 {
		port = 8080 // default fallback
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server running on port %d", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", "error", err)
	}

	logger.Info("Server exiting")
}

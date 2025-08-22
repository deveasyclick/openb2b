package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deveasyclick/openb2b/internal/config"
	"github.com/deveasyclick/openb2b/internal/db"
	"github.com/deveasyclick/openb2b/internal/routes"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/pkg/logger"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn := db.New(db.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		Env:      cfg.Env,
	})

	appCtx := &deps.AppContext{
		DB:     dbConn,
		Config: cfg,
		Logger: logger.New(),
		Cache:  nil,
	}

	routes.Register(r, appCtx)

	srv := http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Printf("Server running on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

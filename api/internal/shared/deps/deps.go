// Package deps provides a centralized application context that bundles
// common dependencies such as database connections, configuration,
// logging, and caching. It is designed to simplify dependency injection
// and reduce boilerplate across the application.
package deps

import (
	"github.com/deveasyclick/openb2b/internal/config"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

// AppContext holds shared application dependencies such as the database,
// configuration, logger, and cache. It is typically passed around to
// different layers of the application to ensure consistent access
// to core services.
type AppContext struct {
	// DB is the GORM database connection used for persistence.
	DB *gorm.DB

	// Config contains application-wide configuration values.
	Config *config.Config

	// Logger provides structured logging capabilities.
	Logger interfaces.Logger

	// Cache provides access to a caching backend (e.g., Redis).
	Cache interfaces.Cache
}

// NewAppContext creates and returns a new AppContext instance with the
// provided dependencies. This function should typically be called once
// during application startup to initialize the shared context.
//
// Example:
//
//	ctx := deps.NewAppContext(db, cfg, logger, cache)
//	service := myservice.NewService(ctx)
func NewAppContext(db *gorm.DB, cfg *config.Config, logger interfaces.Logger, cache interfaces.Cache) *AppContext {
	return &AppContext{
		DB:     db,
		Config: cfg,
		Logger: logger,
		Cache:  cache,
	}
}

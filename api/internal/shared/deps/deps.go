package deps

import (
	"github.com/deveasyclick/openb2b/internal/config"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm"
)

type AppContext struct {
	DB     *gorm.DB
	Config *config.Config
	Logger interfaces.Logger
	Cache  interfaces.Cache
}

func NewAppContext(db *gorm.DB, cfg *config.Config, logger interfaces.Logger, cache interfaces.Cache) *AppContext {
	return &AppContext{
		DB:     db,
		Config: cfg,
		Logger: logger,
		Cache:  cache,
	}
}

package setup

import (
	"net/http/httptest"
	"os"

	"github.com/deveasyclick/openb2b/internal/config"
	"github.com/deveasyclick/openb2b/internal/routes"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/pkg/logger"
	"github.com/go-chi/chi"
)

func SetupTestServer() *httptest.Server {
	r := chi.NewRouter()
	db := SetupTestDB()
	config := &config.Config{
		Env: "test",
	}

	appCtx := &deps.AppContext{
		DB:     db,
		Config: config,                       // or a test config
		Logger: logger.New(os.Getenv("ENV")), // you can use a no-op logger
		Cache:  nil,
	}
	middlewares := NewFake(1, 2, "clerk-user-1")
	routes.Register(r, appCtx, middlewares)
	return httptest.NewServer(r)
}

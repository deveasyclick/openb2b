package db

import (
	"time"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	gormlogger "github.com/deveasyclick/openb2b/pkg/logger/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service struct {
	DB *gorm.DB
}

func New(url string, appLogger interfaces.Logger) *gorm.DB {
	gormLogger := gormlogger.New(appLogger, logger.Info)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		appLogger.Fatal("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Fatal("failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Migrate the schema
	err = db.AutoMigrate(
		&model.Customer{},
		&model.Variant{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.User{},
		&model.Org{},
		&model.Invoice{},
		&model.InvoiceItem{},
	)

	if err != nil {
		appLogger.Fatal("failed to migrate database: %v", err)
	}

	appLogger.Info("Connected to database")

	return db
}

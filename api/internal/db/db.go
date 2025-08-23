package db

import (
	"fmt"
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

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Env      string
}

func New(c DBConfig, appLogger interfaces.Logger) *gorm.DB {
	sslmode := "disable"
	if c.Env == "production" {
		sslmode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		sslmode,
	)
	gormLogger := gormlogger.New(appLogger, logger.Info)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.User{},
		&model.Org{},
		&model.Invoice{},
	)

	if err != nil {
		appLogger.Fatal("failed to migrate database: %v", err)
	}

	appLogger.Info("Connected to database")

	return db
}

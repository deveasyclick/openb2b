package setup

import (
	"log"

	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect test db: %v", err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Org{},
		&model.Product{},
		&model.Variant{},
		&model.Customer{},
		&model.Order{},
		&model.OrderItem{},
	)

	if err != nil {
		log.Fatalf("failed to migrate test database: %v", err)
	}

	TestDB = db

	return db
}

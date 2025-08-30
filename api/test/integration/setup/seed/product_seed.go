package seed

import (
	"log"

	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/gorm"
)

func Insert(db *gorm.DB) {
	db.Create(&model.Product{Name: "Test Product 1"})
	db.Create(&model.Product{Name: "Test Product 2"})
}

// Clear deletes all products and their variants in the database.
func Clear(db *gorm.DB) {
	// Delete variants first (foreign key dependency)
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Variant{}).Error; err != nil {
		log.Fatalf("failed to clear variants: %v", err)
	}

	// Delete products
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Product{}).Error; err != nil {
		log.Fatalf("failed to clear products: %v", err)
	}
}

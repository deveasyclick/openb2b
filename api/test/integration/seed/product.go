package seed

import (
	"log"

	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/gorm"
)

func InsertProducts(db *gorm.DB) model.Product {
	variant := model.Variant{
		SKU:   "SKU-001",
		Color: "Red",
		Size:  "M",
		Price: 10.0,
		Stock: 100,
		BaseModel: model.BaseModel{
			ID: 100,
		},
	}

	variant2 := model.Variant{
		SKU:   "SKU-002",
		Color: "Blue",
		Size:  "L",
		Price: 20.0,
		Stock: 50,
		BaseModel: model.BaseModel{
			ID: 99,
		},
	}

	product := model.Product{
		Name:     "Test Product 1",
		Variants: []model.Variant{variant, variant2},
	}

	err := db.Create(&product).Error
	if err != nil {
		log.Fatalf("failed to create product: %v", err)
	}

	return product
}

// Clear deletes all products and their variants in the database.
func ClearProducts(db *gorm.DB) {
	// Delete variants first (foreign key dependency)
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Product{}).Error; err != nil {
		log.Fatalf("failed to clear orders: %v", err)
	}

}

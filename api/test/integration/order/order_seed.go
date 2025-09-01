package order_test

import (
	"log"

	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/gorm"
)

var nonPendingOrderId uint

func Insert(db *gorm.DB) {
	db.Create(&model.Order{
		Notes:       "Notes",
		OrderNumber: "ORD-123",
	}) // create an order to seed the database
	nonPendingOrder := &model.Order{Status: "processing", OrderNumber: "ORD-124"}
	err := db.Create(nonPendingOrder).Error
	if err != nil {
		log.Fatalf("failed to create order: %v", err)
	}
	nonPendingOrderId = nonPendingOrder.ID
}

// Clear deletes all products and their variants in the database.
func Clear(db *gorm.DB) {
	// Delete variants first (foreign key dependency)
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Order{}).Error; err != nil {
		log.Fatalf("failed to clear orders: %v", err)
	}

}

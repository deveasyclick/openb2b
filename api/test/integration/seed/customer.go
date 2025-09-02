package seed

import (
	"log"

	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/gorm"
)

var Customer model.Customer

func InsertCustomers(db *gorm.DB) {
	Customer = model.Customer{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "i6oP7@example.com",
		PhoneNumber: "+1-202-555-0199",
		Address: &model.Address{
			Address: "123 Market Street",
			City:    "San Francisco",
			State:   "California",
			Country: "USA",
			Zip:     "02912",
		},
		Company: "OpenB2B",
	}

	err := db.Create(&Customer).Error
	if err != nil {
		log.Fatalf("failed to create customer: %v", err)
	}

}

func ClearCustomers(db *gorm.DB) {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Customer{}).Error; err != nil {
		log.Fatalf("failed to clear orders: %v", err)
	}

}

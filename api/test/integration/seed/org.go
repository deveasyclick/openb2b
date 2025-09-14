package seed

import (
	"log"

	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/gorm"
)

var Org model.Org

func InsertOrgs(db *gorm.DB) {
	Org = model.Org{
		Name:             "OpenB2B",
		Logo:             "https://www.openb2b.com/logo.png",
		OrganizationName: "OpenB2B NG",
		OrganizationUrl:  "https://www.openb2b.com",
		Email:            "info@openb2b.com",
		Phone:            "+1-202-555-0199",
		Address: &model.Address{
			Address: "123 Market Street",
			City:    "San Francisco",
			State:   "California",
			Country: "USA",
			Zip:     "02912",
		},
	}

	err := db.Create(&Org).Error
	if err != nil {
		log.Fatalf("failed to create org: %v", err)
	}

}

func ClearOrgs(db *gorm.DB) {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Org{}).Error; err != nil {
		log.Fatalf("failed to clear orders: %v", err)
	}

}

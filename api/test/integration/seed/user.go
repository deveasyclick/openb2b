package seed

import (
	"log"

	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/gorm"
)

func InsertUsers(db *gorm.DB) {
	user := model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "i6oP7@example.com",
		Role:      model.RoleAdmin,
		Address: &model.Address{
			Address: "123 Market Street",
			City:    "San Francisco",
			State:   "California",
			Country: "USA",
			Zip:     "02912",
		},
	}

	err := db.Create(&user).Error
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

}

func ClearUsers(db *gorm.DB) {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.User{}).Error; err != nil {
		log.Fatalf("failed to clear users: %v", err)
	}

}

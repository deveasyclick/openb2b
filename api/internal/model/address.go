package model

type Address struct {
	Zip     string `gorm:"type:varchar(100)" json:"zip"`
	State   string `gorm:"type:varchar(50)" json:"state"`
	City    string `gorm:"type:varchar(50)" json:"city"`
	Country string `gorm:"type:varchar(50)" json:"country"`
	Address string `gorm:"type:varchar(200)" json:"address"`
}

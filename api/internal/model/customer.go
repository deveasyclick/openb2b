package model

// Customer represents a customer belonging to a specific org.
// To ensure a customer is unique within a org, we enforce a composite unique index
// on (org_id, phone_number). Email is not used for uniqueness because it is optional.
type Customer struct {
	BaseModel
	FirstName   string   `gorm:"not null;type:varchar(100);check:first_name <> ''" json:"firstName" validate:"required,max=100"`
	LastName    string   `gorm:"not null;type:varchar(100);check:last_name <> ''" json:"lastName" validate:"required,max=100"`
	PhoneNumber string   `gorm:"index:uniqueIndex:idx_org_phone;not null;type:varchar(50);check:phone_number <> ''" json:"phoneNumber"`
	Email       string   `gorm:"index;type:varchar(100);" json:"email"`
	Address     *Address `gorm:"embedded;embeddedPrefix:address_" json:"address"`
	Company     string   `gorm:"type:varchar(100)" json:"company"`
	OrgID       uint     `gorm:"index" json:"orgId"`
	Org         *Org     `gorm:"foreignKey:OrgID" json:"org,omitempty"`
	Orders      []*Order `json:"orders,omitempty"`
}

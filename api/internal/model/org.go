package model

// Org represents an organization entity
// @Description Organization response model
type Org struct {
	BaseModel
	Name             string     `gorm:"not null;unique;index;type:varchar(50);check:name <> ''" json:"name" validate:"required,max=50"`
	Logo             string     `gorm:"type:varchar(100)" json:"logo"`
	OrganizationName string     `gorm:"not null;type:varchar(50);check:organization_name <> ''" json:"organizationName" validate:"required,max=50"`
	OrganizationUrl  string     `gorm:"type:varchar(100);" json:"organizationUrl" validate:"max=100"`
	Email            string     `gorm:"unique;type:varchar(50);check:email <> ''" json:"email" validate:"required,max=50"`
	Phone            string     `gorm:"not null;type:varchar(50);check:phone <> ''" json:"phone" validate:"required,max=50"`
	Address          *Address   `gorm:"embedded;embeddedPrefix:address_" json:"address"`
	OnboardedAt      bool       `gorm:"type:boolean;default:false" json:"onboardedAt"`
	Users            []Customer `json:"users,omitempty"`
	Products         []Product  `json:"products,omitempty"`
	Customers        []Customer `json:"customers,omitempty"`
	Orders           []Order    `json:"orders,omitempty"`
}

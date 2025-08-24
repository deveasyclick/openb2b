package model

import (
	"database/sql/driver"
	"errors"
)

type Role string

const (
	RoleOwner  Role = "distributor"
	RoleAdmin  Role = "admin"
	RoleViewer Role = "viewer"
)

const errInvalidRoleValue = "invalid role value"

// Scan implements the Scanner interface for database deserialization
func (r *Role) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New(errInvalidRoleValue)
	}

	switch Role(strValue) {
	case RoleOwner, RoleAdmin, RoleViewer:
		*r = Role(strValue)
		return nil
	default:
		return errors.New(errInvalidRoleValue)
	}
}

// Value implements the driver Valuer interface for database serialization
func (r Role) Value() (driver.Value, error) {
	switch r {
	case RoleOwner, RoleAdmin, RoleViewer:
		return string(r), nil
	default:
		return nil, errors.New(errInvalidRoleValue)
	}
}

type User struct {
	BaseModel
	ClerkID   string  `gorm:"uniqueIndex;type:varchar(50)" json:"clerkId"`
	FirstName string  `gorm:"not null;type:varchar(100);check:first_name <> ''" json:"firstName" validate:"required,max=100"`
	LastName  string  `gorm:"not null;type:varchar(100);check:last_name <> ''" json:"lastName" validate:"required,max=100"`
	Email     string  `gorm:"uniqueIndex;type:varchar(50);check:email <> ''" json:"email" validate:"required,max=50"`
	Phone     *string `gorm:"type:varchar(50)" json:"phone" validate:"omitempty,max=50"`
	Role      string  `gorm:"type:enum('owner','admin','sales','viewer');default:'sales'" json:"role"`
	OrgID     uint    `gorm:"index" json:"orgId"`
	Org       Org     `gorm:"foreignKey:OrgID" json:"org,omitempty"`
	Address   Address `gorm:"embedded;embeddedPrefix:address_" json:"address"`
}

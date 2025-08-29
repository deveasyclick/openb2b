package model

// Product represents an product entity
// @Description Product response model
type Product struct {
	BaseModel
	Name        string    `gorm:"not null" json:"name"`
	Category    string    `json:"category"`
	OrgID       uint      `gorm:"index;not null" json:"orgId"`
	Org         *Org      `gorm:"foreignKey:OrgID" json:"org,omitempty"`
	ImageURL    string    `json:"imageUrl"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants" gorm:"foreignKey:ProductID"`
}

// Variant represents an variant entity
// @Description Variant response model
type Variant struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint `gorm:"index;not null"`
	Color     string
	Size      string
	Price     float64 `gorm:"not null"`
	Stock     int     `gorm:"not null"`
	SKU       string  `gorm:"not null;uniqueIndex:idx_org_sku"`
	OrgID     uint    `gorm:"not null;uniqueIndex:idx_org_sku"` //needed for sku uniqueness per org
}

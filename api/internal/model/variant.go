package model

// Variant represents an variant entity
// @Description Variant response model
type Variant struct {
	BaseModel
	ProductID uint    `gorm:"index;not null" json:"productId"`
	Color     string  `json:"color"`
	Size      string  `json:"size"`
	Price     float64 `gorm:"not null" json:"price"`
	Stock     int     `gorm:"not null" json:"stock"`
	SKU       string  `gorm:"not null;uniqueIndex:idx_org_sku" json:"sku"`
	OrgID     uint    `gorm:"not null;uniqueIndex:idx_org_sku" json:"orgId"` //needed for sku uniqueness per org
}

package model

type Product struct {
	BaseModel
	Name     string  `gorm:"not null" json:"name"`
	SKU      string  `gorm:"not null;uniqueIndex:idx_org_sku" json:"sku"`
	Category string  `json:"category"`
	Price    float64 `gorm:"not null" json:"price"`
	Quantity int     `json:"quantity"`
	OrgID    uint    `gorm:"index;uniqueIndex:idx_org_sku" json:"orgId"`
	Org      Org     `gorm:"foreignKey:OrgID" json:"org,omitempty"`
}

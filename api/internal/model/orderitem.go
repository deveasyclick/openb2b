package model

type OrderItem struct {
	BaseModel

	OrderID   uint    `json:"orderId"`
	ProductID uint    `json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	VariantID uint    `json:"variantId"`
	Variant   Variant `gorm:"foreignKey:VariantID" json:"variant"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Total     float64 `json:"total"`
	OrgID     uint    `json:"orgId"`
}

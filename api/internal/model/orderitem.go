package model

type OrderItem struct {
	BaseModel

	OrderID   uint     `json:"orderId" gorm:"uniqueIndex:idx_order_variant"`
	ProductID uint     `json:"productId"`
	Product   *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	VariantID uint     `gorm:"uniqueIndex:idx_order_variant" json:"variantId"`
	Variant   *Variant `gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Total     float64 `json:"total"` // (UnitPrice*Qty - discounts) + tax
	OrgID     uint    `json:"orgId"`

	TaxRate   float64 `json:"taxRate"`   // e.g., 0.10 for 10%
	TaxAmount float64 `json:"taxAmount"` // tax charged on this line (after discounts)
	Notes     string  `json:"notes"`

	Discount             DiscountInfo `gorm:"embedded;embeddedPrefix:discount_" json:"discount"`
	AppliedDiscount      float64      `json:"appliedDiscount"`      // Actual discount applied
	AppliedOrderDiscount float64      `json:"appliedOrderDiscount"` // proportional share of order-level discount
}

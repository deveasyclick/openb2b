package model

type InvoiceItem struct {
	BaseModel

	OrgID uint `gorm:"not null;index"`
	Org   *Org `gorm:"foreignKey:OrgID" json:"org"`

	InvoiceID uint    `gorm:"not null;index"` // Link to invoice
	VariantID uint    `gorm:"not null;index"`
	Variant   Variant `gorm:"foreignKey:VariantID" json:"variant"`

	Notes     string  `gorm:"type:varchar(255)" json:"description"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
	UnitPrice float64 `gorm:"not null;default:0" json:"unitPrice"`
	TaxAmount float64 `gorm:"not null;default:0" json:"taxAmount"`
	LineTotal float64 `gorm:"not null;default:0" json:"lineTotal"` // Quantity * UnitPrice + TaxAmount
	Subtotal  float64 `json:"subtotal"`                            // Sum of all item totals before discount and tax
}

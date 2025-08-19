package model

type Invoice struct {
	BaseModel
	OrderID       uint   `gorm:"index;not null" json:"orderId"`
	Order         Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	InvoiceNumber string `gorm:"unique;not null" json:"invoiceNumber"`
	PDFPath       string `json:"pdfPath"`
	OrgID         uint   `gorm:"index" json:"orgId"`
	Org           Org    `gorm:"foreignKey:OrgID" json:"org,omitempty"`
}

package model

import "time"

// InvoiceStatus represents possible statuses for invoices
type InvoiceStatus string

const (
	//draft, Internal, editable, not sent to customer.
	InvoiceStatusDraft InvoiceStatus = "draft"
	// pro_forma, Sent to customer for review (not a legal invoice).
	InvoiceStatusProForma InvoiceStatus = "pro_forma"
	// issued, Official invoice, data locked, triggers accounting.
	InvoiceStatusIssued InvoiceStatus = "issued"
	// paid, Fully paid, invoice closed.
	InvoiceStatusPaid InvoiceStatus = "paid"
	// overdue, Invoice overdue, customer has not paid.
	InvoiceStatusOverdue InvoiceStatus = "overdue"
	// cancelled, Invoice cancelled, customer has not paid.
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
	// partially_paid, Invoice partially paid, customer has not paid.
	InvoiceStatusPartiallyPaid InvoiceStatus = "partially_paid"
)

// Invoice represents an invoice linked to an order
type Invoice struct {
	BaseModel

	OrgID   uint   `gorm:"index;not null" json:"orgId"`
	OrderID uint   `gorm:"index;not null" json:"orderId"`
	Order   *Order `gorm:"foreignKey:OrderID" json:"order"`

	InvoiceNumber string        `gorm:"uniqueIndex;size:50;not null" json:"invoiceNumber"`
	Status        InvoiceStatus `gorm:"type:varchar(20);default:'draft';not null" json:"status"`

	IssuedAt time.Time  `gorm:"not null" json:"issuedAt"`
	DueDate  *time.Time `json:"dueDate"`

	Currency      string  `gorm:"size:3;default:'NGN';not null" json:"currency"`
	Subtotal      float64 `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	TaxTotal      float64 `gorm:"type:decimal(12,2);not null" json:"taxTotal"`
	DiscountTotal float64 `gorm:"type:decimal(12,2);not null" json:"discountTotal"`
	Total         float64 `gorm:"type:decimal(12,2);not null" json:"total"`

	Notes  string `gorm:"type:text" json:"notes"`
	PDFUrl string `gorm:"type:text" json:"pdf_url"`

	Items []*InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items"`
}

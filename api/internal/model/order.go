package model

import (
	"time"

	"gorm.io/gorm"
)

type DeliveryStatus string
type OrderStatus string
type DiscountType string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusApproved  OrderStatus = "approved"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"

	DeliveryPending   DeliveryStatus = "pending"
	DeliveryShipped   DeliveryStatus = "shipped"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryCancelled DeliveryStatus = "cancelled"

	DiscountPercentage DiscountType = "percentage"
	DiscountFixed      DiscountType = "fixed"
)

type DeliveryInfo struct {
	Address       *Address       `gorm:"embedded;embeddedPrefix:address_" json:"address"`
	TransportFare float64        `gorm:"not null" json:"transportFare"`
	Status        DeliveryStatus `gorm:"type:varchar(20)" json:"status"`
	Date          *time.Time     `json:"date"`
	At            *time.Time     `json:"at"`
}

type DiscountInfo struct {
	Type   DiscountType `gorm:"not null" json:"type"`
	Amount float64      `gorm:"not null" json:"amount"`
}

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	if o.Delivery.Status == DeliveryDelivered && o.Delivery.At == nil {
		now := time.Now()
		o.Delivery.At = &now
	}
	return
}

// Add customer instead of collecting customer name and phone number to prevent redundancy, preserve customer order history and stats
type Order struct {
	BaseModel

	OrderNumber string       `gorm:"uniqueIndex;size:50" json:"orderNumber"`
	CustomerID  uint         `gorm:"index;not null" json:"customerId"`
	Customer    *Customer    `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Status      OrderStatus  `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','processing','shipped','completed')" json:"status"`
	OrgID       uint         `gorm:"index" json:"orgId"`
	Org         *Org         `gorm:"foreignKey:OrgID" json:"org"`
	Items       []OrderItem  `gorm:"foreignKey:OrderID" json:"items"`
	Delivery    DeliveryInfo `gorm:"embedded;embeddedPrefix:delivery_" json:"delivery"`
	Notes       string       `json:"notes"`

	Discount          DiscountInfo `gorm:"embedded;embeddedPrefix:discount_" json:"discount"`
	AppliedDiscount   float64      `json:"appliedDiscount"` // Actual discount applied
	DiscountTotal     float64      `json:"discountTotal"`   /// ItemDiscountTotal + AppliedDiscount
	ItemDiscountTotal float64      // sum of all per-item discounts

	Total    float64 `json:"total"`     // final payable amount = sum of all item totals
	Subtotal float64 `json:"subtotal"`  // sum of item (unitPrice * qty), before discounts & tax
	TaxTotal float64 `json:"taxAmount"` // Sum of all item tax amounts

	Invoices []Invoice `gorm:"foreignKey:OrderID" json:"invoices"`
}

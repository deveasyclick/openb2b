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
	Address       string         `gorm:"not null" json:"address"`
	City          string         `gorm:"not null;type:varchar(100)" json:"city"`
	State         string         `gorm:"not null;type:varchar(100)" json:"state"`
	Country       string         `gorm:"not null;type:varchar(100)" json:"country"`
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
	DistributorID uint         `gorm:"index" json:"distributorId"`
	CreatedBy     *User        `gorm:"foreignKey:DistributorID" json:"createdBy"`
	CustomerID    uint         `gorm:"index" json:"customerId"`
	Customer      *Customer    `gorm:"foreignKey:CustomerID" json:"requestedFor"`
	Status        OrderStatus  `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','processing','shipped','completed')" json:"status"`
	OrgID         uint         `gorm:"index" json:"orgId"`
	Org           *Org         `gorm:"foreignKey:OrgID" json:"org"`
	OrderItems    []OrderItem  `gorm:"foreignKey:OrderID" json:"orderItems"`
	Delivery      DeliveryInfo `gorm:"embedded;embeddedPrefix:delivery_" json:"delivery"`
	Notes         string       `json:"notes"`
	Discount      DiscountInfo `gorm:"embedded;embeddedPrefix:discount_" json:"discount"`
}

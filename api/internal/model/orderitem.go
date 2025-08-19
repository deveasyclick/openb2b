package model

type OrderItem struct {
	BaseModel
	OrderID   uint    `gorm:"index" json:"orderId"`
	Order     Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductID uint    `gorm:"index" json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"`
}

package order

import (
	"time"

	"github.com/deveasyclick/openb2b/internal/model"
)

//
// CREATE DTOs
//

type CreateOrderItemDTO struct {
	VariantID uint    `json:"variantId" validate:"required"`
	ProductID uint    `json:"productId" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	UnitPrice float64 `json:"unitPrice" validate:"required,gt=0"`
}

func (i *CreateOrderItemDTO) ToModel(orgID uint) model.OrderItem {
	item := model.OrderItem{
		VariantID: i.VariantID,
		Quantity:  i.Quantity,
		UnitPrice: i.UnitPrice,
		Total:     float64(i.Quantity) * i.UnitPrice,
		OrgID:     orgID,
		ProductID: i.ProductID,
	}

	return item
}

type CreateDeliveryInfoDTO struct {
	Address       *model.Address `json:"address" validate:"required"`
	TransportFare float64        `json:"transportFare" validate:"required,min=0"`
}

func (d *CreateDeliveryInfoDTO) ToModel() model.DeliveryInfo {
	return model.DeliveryInfo{
		Address:       d.Address,
		TransportFare: d.TransportFare,
	}
}

type CreateDiscountInfoDTO struct {
	Type   model.DiscountType `json:"type" validate:"required,oneof=percentage fixed"`
	Amount float64            `json:"amount" validate:"required,min=0"`
}

func (di *CreateDiscountInfoDTO) ToModel() model.DiscountInfo {
	return model.DiscountInfo{
		Type:   di.Type,
		Amount: di.Amount,
	}
}

type CreateOrderDTO struct {
	CustomerID uint                  `json:"customerId" validate:"required"`
	Items      []CreateOrderItemDTO  `json:"items" validate:"required,dive"`
	Delivery   CreateDeliveryInfoDTO `json:"delivery" validate:"required"`
	Notes      string                `json:"notes" validate:"omitempty,max=1000"`
	Discount   CreateDiscountInfoDTO `json:"discount" validate:"omitempty"`
	Tax        float64               `json:"tax" validate:"min=0"`
}

func (dto *CreateOrderDTO) ToModel(orgID uint) model.Order {
	order := model.Order{
		OrderNumber: generateOrderNumber(),
		CustomerID:  dto.CustomerID,
		OrgID:       orgID,
		Delivery:    dto.Delivery.ToModel(),
		Notes:       dto.Notes,
		Discount:    dto.Discount.ToModel(),
		Tax:         dto.Tax,
		Status:      model.OrderStatusPending,
	}
	// Map items
	for _, item := range dto.Items {
		order.Items = append(order.Items, item.ToModel(orgID))
	}

	calculateTotals(&order)

	return order
}

//
// UPDATE DTOs
//

type UpdateOrderDTO struct {
	Status   *model.OrderStatus     `json:"status" validate:"omitempty,oneof=pending approved delivered cancelled"`
	Notes    *string                `json:"notes" validate:"omitempty,max=1000"`
	Discount *CreateDiscountInfoDTO `json:"discount" validate:"omitempty"`
	Tax      *float64               `json:"tax" validate:"omitempty,min=0"`
	Items    []*CreateOrderItemDTO  `json:"items" validate:"omitempty,dive"`
	Delivery *UpdateDeliveryInfoDTO `json:"deliver" validate:"omitempty"`
}

func (dto *UpdateOrderDTO) ApplyModel(order *model.Order) {
	if dto.Status != nil {
		order.Status = *dto.Status
	}
	if dto.Notes != nil {
		order.Notes = *dto.Notes
	}
	if dto.Discount != nil {
		order.Discount = dto.Discount.ToModel()
	}
	if dto.Tax != nil {
		order.Tax = *dto.Tax
	}

	if dto.Delivery != nil {
		dto.Delivery.ApplyModel(&order.Delivery)
	}

	if dto.Items != nil {
		var updatedItems []model.OrderItem
		for _, updatedItem := range dto.Items {
			updatedItems = append(order.Items, updatedItem.ToModel(order.OrgID))
		}
		order.Items = updatedItems
	}

	calculateTotals(order)
}

type UpdateDeliveryInfoDTO struct {
	Address       *model.Address        `json:"address" validate:"omitempty"`
	TransportFare *float64              `json:"transportFare" validate:"omitempty,min=0"`
	Status        *model.DeliveryStatus `json:"status" validate:"omitempty,oneof=pending shipped delivered cancelled"`
	Date          *time.Time            `json:"date"`
}

func (dto *UpdateDeliveryInfoDTO) ApplyModel(delivery *model.DeliveryInfo) {
	if dto.Address != nil {
		delivery.Address = dto.Address
	}
	if dto.TransportFare != nil {
		delivery.TransportFare = *dto.TransportFare
	}
	if dto.Status != nil {
		delivery.Status = *dto.Status
	}
	if dto.Date != nil {
		delivery.Date = dto.Date
	}
}

type UpdateOrderItemDTO struct {
	Quantity  *int     `json:"quantity" validate:"omitempty,min=1"`
	UnitPrice *float64 `json:"unitPrice" validate:"omitempty,gt=0"`
}

func (dto *UpdateOrderItemDTO) ApplyModel(item *model.OrderItem) {
	if dto.Quantity != nil {
		item.Quantity = *dto.Quantity
	}
	if dto.UnitPrice != nil {
		item.UnitPrice = *dto.UnitPrice
	}
	item.Total = float64(item.Quantity) * item.UnitPrice
}

// Calculate totals
func calculateTotals(order *model.Order) {
	// Calculate subtotal
	subtotal := 0.0
	for _, item := range order.Items {
		subtotal += item.Total
	}
	order.Subtotal = subtotal

	// Calculate discount
	discountAmount := 0.0
	discount := order.Discount
	if discount.Type == "percentage" {
		discountAmount = subtotal * (discount.Amount / 100)
	} else if discount.Type == "fixed" {
		discountAmount = discount.Amount
	}
	order.DiscountAmount = discountAmount

	// Calculate tax
	taxAmount := 0.0
	if order.Tax > 0 {
		taxAmount = (subtotal - discountAmount) * (order.Tax / 100)
	}

	// Final total
	order.Total = subtotal - discountAmount + taxAmount
	if order.Total < 0 {
		order.Total = 0
	}
}

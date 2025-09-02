package dto

import (
	"time"

	"github.com/deveasyclick/openb2b/internal/model"
	generateordernumber "github.com/deveasyclick/openb2b/internal/utils/generateOrderNumber"
)

//
// CREATE DTOs
//

type CreateOrderItemDTO struct {
	VariantID uint `json:"variantId" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

func (i *CreateOrderItemDTO) ToModel(orgID uint, variant model.Variant) model.OrderItem {
	item := model.OrderItem{
		VariantID: i.VariantID,
		Quantity:  i.Quantity,
		UnitPrice: variant.Price,
		Total:     float64(i.Quantity) * variant.Price,
		OrgID:     orgID,
		ProductID: variant.ProductID,
	}

	return item
}

type CreateDeliveryInfoDTO struct {
	Address       *model.Address `json:"address" validate:"required"`
	TransportFare float64        `json:"transportFare" validate:"min=0"`
}

func (d *CreateDeliveryInfoDTO) ToModel() model.DeliveryInfo {
	return model.DeliveryInfo{
		Address:       d.Address,
		TransportFare: d.TransportFare,
	}
}

type CreateDiscountInfoDTO struct {
	Type   model.DiscountType `json:"type" validate:"required,oneof=percentage fixed"`
	Amount float64            `json:"amount" validate:"min=0"`
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

func (dto *CreateOrderDTO) ToModel(variantMap map[uint]model.Variant, orgID uint) model.Order {
	order := model.Order{
		OrderNumber: generateordernumber.GenerateOrderNumber(),
		CustomerID:  dto.CustomerID,
		OrgID:       orgID,
		Delivery:    dto.Delivery.ToModel(),
		Notes:       dto.Notes,
		Discount:    dto.Discount.ToModel(),
		Tax:         dto.Tax,
		Status:      model.OrderStatusPending,
	}

	for _, item := range dto.Items {
		v := variantMap[item.VariantID]
		order.Items = append(order.Items, model.OrderItem{
			VariantID: v.ID,
			ProductID: v.ProductID,
			UnitPrice: v.Price,
			Quantity:  item.Quantity,
			Total:     float64(item.Quantity) * v.Price,
			OrgID:     orgID,
		})
	}

	calculateTotals(&order)
	return order
}

//
// UPDATE DTOs
//

type UpdateOrderDTO struct {
	Status     *model.OrderStatus     `json:"status" validate:"omitempty,oneof=pending approved delivered cancelled"`
	Notes      *string                `json:"notes" validate:"omitempty,max=1000"`
	Discount   *CreateDiscountInfoDTO `json:"discount" validate:"omitempty"`
	Tax        *float64               `json:"tax" validate:"omitempty,min=0"`
	Items      []*CreateOrderItemDTO  `json:"items" validate:"omitempty,dive"`
	Delivery   *UpdateDeliveryInfoDTO `json:"deliver" validate:"omitempty"`
	CustomerID *uint                  `json:"customerId" validate:"omitempty"`
}

func (dto *UpdateOrderDTO) ApplyModel(order *model.Order, variantMap *map[uint]model.Variant) {
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

	if dto.CustomerID != nil {
		order.CustomerID = *dto.CustomerID
	}

	if dto.Items != nil && variantMap != nil {
		for _, item := range dto.Items {
			v := (*variantMap)[item.VariantID]
			order.Items = append(order.Items, model.OrderItem{
				VariantID: v.ID,
				ProductID: v.ProductID,
				UnitPrice: v.Price,
				Quantity:  item.Quantity,
				Total:     float64(item.Quantity) * v.Price,
				OrgID:     order.OrgID,
			})
		}

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
	order.TaxAmount = taxAmount

	// Final total
	order.Total = subtotal - discountAmount + taxAmount
	if order.Total < 0 {
		order.Total = 0
	}
}

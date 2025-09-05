package dto

import (
	"time"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/utils/numbergen"
	"github.com/deveasyclick/openb2b/internal/utils/ordertotals"
)

// CreateOrderItemDTO represents incoming data for creating an order item
type CreateOrderItemDTO struct {
	VariantID uint                  `json:"variantId" validate:"required"`
	Quantity  int                   `json:"quantity" validate:"required,min=1"`
	Discount  CreateDiscountInfoDTO `json:"discount" validate:"omitempty"`
	Notes     string                `json:"notes" validate:"omitempty"`
}

// ToModel converts CreateOrderItemDTO to a fully initialized OrderItem
// Order items won'te be created separately, they are created when the order is created so we don't need to calculate totals at item level
func (i *CreateOrderItemDTO) ToModel(orgID uint, variant model.Variant) model.OrderItem {

	return model.OrderItem{
		OrgID:     orgID,
		ProductID: variant.ProductID,
		VariantID: i.VariantID,
		Notes:     i.Notes,
		Quantity:  i.Quantity,
		UnitPrice: variant.Price,
		Discount: model.DiscountInfo{
			Type:   i.Discount.Type,
			Amount: i.Discount.Amount,
		},
		TaxRate: variant.TaxRate,
	}
}

type CreateDeliveryInfoDTO struct {
	Address       AddressRequired `json:"address" validate:"required"`
	TransportFare float64         `json:"transportFare" validate:"min=0"`
}

func (d *CreateDeliveryInfoDTO) ToModel() model.DeliveryInfo {
	return model.DeliveryInfo{
		Address:       d.Address.ToModel(),
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
}

func (dto *CreateOrderDTO) ToModel(variantMap map[uint]model.Variant, orgID uint) model.Order {
	order := model.Order{
		OrderNumber: numbergen.Generate("ORD"),
		CustomerID:  dto.CustomerID,
		OrgID:       orgID,
		Delivery:    dto.Delivery.ToModel(),
		Notes:       dto.Notes,
		Discount:    dto.Discount.ToModel(),
		Status:      model.OrderStatusPending,
		Items:       make([]model.OrderItem, 0, len(dto.Items)),
	}

	// Convert each DTO item to OrderItem model
	for _, itemDTO := range dto.Items {
		variant, ok := variantMap[itemDTO.VariantID]
		if !ok {
			continue // skip invalid variants
		}

		orderItem := itemDTO.ToModel(orgID, variant)
		order.Items = append(order.Items, orderItem)
	}

	// After items are added, calculate totals including applied order-level discount
	ordertotals.Calculate(&order)

	return order
}

//
// UPDATE DTOs
//

type UpdateOrderDTO struct {
	Status     *model.OrderStatus     `json:"status" validate:"omitempty,oneof=pending approved delivered cancelled"`
	Notes      *string                `json:"notes" validate:"omitempty,max=1000"`
	Discount   *CreateDiscountInfoDTO `json:"discount" validate:"omitempty"`
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

	if dto.Delivery != nil {
		dto.Delivery.ApplyModel(&order.Delivery)
	}

	if dto.CustomerID != nil {
		order.CustomerID = *dto.CustomerID
	}

	if len(dto.Items) > 0 && len(*variantMap) > 0 {
		order.Items = make([]model.OrderItem, 0, len(dto.Items))
		for _, item := range dto.Items {
			v, ok := (*variantMap)[item.VariantID]
			if !ok {
				continue // or handle error
			}
			order.Items = append(order.Items, item.ToModel(order.OrgID, v))
		}
	}

	ordertotals.Calculate(order)
}

type UpdateDeliveryInfoDTO struct {
	Address       *model.Address        `json:"address" validate:"omitempty"`
	TransportFare *float64              `json:"transportFare" validate:"omitempty,min=0"`
	Status        *model.DeliveryStatus `json:"status" validate:"omitempty,oneof=pending shipped delivered cancelled"`
	Date          *time.Time            `json:"date" validate:"omitempty,datetime"`
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

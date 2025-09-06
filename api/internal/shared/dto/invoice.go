package dto

import (
	"time"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/utils/numbergen"
)

type CreateInvoiceDTO struct {
	OrderID uint   `json:"orderId" validate:"required"`
	Notes   string `json:"notes,omitempty"`
}

// CreateInvoiceDTO represents incoming API payload to create an invoice
func (d *CreateInvoiceDTO) ToModel(
	orgID uint,
	order *model.Order,
) *model.Invoice {
	inv := &model.Invoice{
		OrgID:           orgID,
		OrderID:         order.ID,
		InvoiceNumber:   numbergen.Generate("INV"),
		Notes:           d.Notes,
		IssuedAt:        time.Now(),
		Status:          model.InvoiceStatusDraft,
		CustomerEmail:   order.Customer.Email,
		CustomerPhone:   order.Customer.PhoneNumber,
		CustomerName:    order.Customer.FirstName + " " + order.Customer.LastName,
		CustomerAddress: order.Customer.Address,

		// Snapshot financial data
		Subtotal:      order.Subtotal,
		TaxTotal:      order.TaxTotal,
		DiscountTotal: order.DiscountTotal,
		Total:         order.Total,
	}

	// Copy order items into invoice items
	inv.Items = make([]*model.InvoiceItem, len(order.Items))
	for i, oi := range order.Items {
		inv.Items[i] = &model.InvoiceItem{
			OrgID:     orgID,
			VariantID: oi.VariantID,
			Notes:     oi.Notes,
			Quantity:  oi.Quantity,
			UnitPrice: oi.UnitPrice,
			TaxAmount: oi.TaxAmount,
			LineTotal: oi.Total,
			Subtotal:  float64(oi.Quantity) * oi.UnitPrice,
			SKU:       oi.SKU,
		}
	}

	return inv
}

// UpdateInvoiceDTO allows updating certain fields (e.g., status, notes)
type UpdateInvoiceDTO struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=draft issued paid partially_paid cancelled"`
	Notes  *string `json:"notes,omitempty"`
	// Allow optional update of DueDate
	DueDate *time.Time `json:"dueDate,omitempty"`
}

// ApplyModel updates allowed fields on an Invoice
func (dto *UpdateInvoiceDTO) ApplyModel(inv *model.Invoice) {
	if dto.Status != nil {
		inv.Status = model.InvoiceStatus(*dto.Status)
	}
	if dto.Notes != nil {
		inv.Notes = *dto.Notes
	}
	if dto.DueDate != nil {
		inv.DueDate = dto.DueDate
	}
}

package order_test

import "github.com/deveasyclick/openb2b/internal/model"

// Calculate totals
func calculateTotals(order model.Order) (subtotal float64, discountAmount float64, taxAmount float64, total float64) {
	// Calculate subtotal
	for _, item := range order.Items {
		subtotal += item.Total
	}

	// Calculate discount
	discount := order.Discount
	if discount.Type == "percentage" {
		discountAmount = subtotal * (discount.Amount / 100)
	} else if discount.Type == "fixed" {
		discountAmount = discount.Amount
	}

	// Calculate tax
	if order.Tax > 0 {
		taxAmount = (subtotal - discountAmount) * (order.Tax / 100)
	}

	// Final total
	total = subtotal - discountAmount + taxAmount
	if total < 0 {
		total = 0
	}

	return subtotal, discountAmount, taxAmount, total
}

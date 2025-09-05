// Package ordertotals provides utilities to calculate discounts, taxes,
// and totals for orders and their items.
package ordertotals

import (
	"math"

	"github.com/deveasyclick/openb2b/internal/model"
)

// round2 rounds a float64 value to 2 decimal places (e.g., cents).
func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

// calculateItemDiscount calculates the discount applied to a single order item.
// The discount is determined by the item's discount type (percentage or fixed amount).
// The discount cannot exceed the item's subtotal.
func calculateItemDiscount(item *model.OrderItem) float64 {
	subtotal := item.UnitPrice * float64(item.Quantity)
	var discount float64

	switch item.Discount.Type {
	case model.DiscountPercentage:
		discount = subtotal * (item.Discount.Amount / 100)
	case model.DiscountFixed:
		discount = item.Discount.Amount
	default:
		discount = 0
	}

	// Ensure discount does not exceed subtotal
	if discount > subtotal {
		discount = subtotal
	}

	return round2(discount)
}

// calculateOrderDiscount calculates the total discount applied at the order level.
// The discount is based on the order's discount type (percentage or fixed amount).
// The maximum discount is capped so the combined item- and order-level discounts
// cannot exceed the subtotal.
func calculateOrderDiscount(order *model.Order) float64 {
	subtotal := order.Subtotal
	itemDiscountTotal := order.ItemDiscountTotal
	var discount float64

	switch order.Discount.Type {
	case model.DiscountPercentage:
		discount = subtotal * (order.Discount.Amount / 100)
	case model.DiscountFixed:
		discount = order.Discount.Amount
	default:
		discount = 0
	}

	// Prevent total discounts from exceeding subtotal
	maxDiscount := subtotal - itemDiscountTotal
	if discount > maxDiscount {
		discount = maxDiscount
	}

	return round2(discount)
}

// applyOrderDiscountToItems distributes the total order-level discount proportionally
// to all order items based on their share of the subtotal.
// The last item gets any rounding adjustment to ensure totals match exactly.
func applyOrderDiscountToItems(order *model.Order) {
	subtotal := order.Subtotal
	if order.AppliedDiscount <= 0 || subtotal <= 0 {
		return
	}

	remaining := order.AppliedDiscount

	for i := range order.Items {
		item := &order.Items[i]
		lineSubtotal := item.UnitPrice * float64(item.Quantity)

		// Proportionally allocate discount
		applied := (lineSubtotal / subtotal) * order.AppliedDiscount

		// Round and adjust remaining for the last item
		if i == len(order.Items)-1 {
			applied = remaining
		} else {
			applied = round2(applied)
			remaining -= applied
		}

		item.AppliedOrderDiscount = applied
	}
}

// calculateTaxAndLineTotal recalculates the tax amount and total line cost for an item.
// Tax is applied on the item's taxable amount (price - discounts).
func calculateTaxAndLineTotal(item *model.OrderItem) {
	// Compute taxable base
	taxable := item.UnitPrice*float64(item.Quantity) - item.AppliedDiscount - item.AppliedOrderDiscount
	if taxable < 0 {
		taxable = 0
	}

	// Calculate tax and total
	item.TaxAmount = round2(taxable * item.TaxRate)
	item.Total = round2(taxable + item.TaxAmount)
}

// Calculate recalculates all financial fields of an order, including:
// - Item-level discounts
// - Order-level discount
// - Tax amounts
// - Final totals
// It updates the order in place.
func Calculate(order *model.Order) {
	// Handle empty orders
	if len(order.Items) == 0 {
		order.Subtotal = 0
		order.ItemDiscountTotal = 0
		order.AppliedDiscount = 0
		order.DiscountTotal = 0
		order.TaxTotal = 0
		order.Total = 0
		return
	}

	var subtotal float64
	var itemDiscountTotal float64

	// Step 1: calculate per-item discounts and subtotal
	for i := range order.Items {
		item := &order.Items[i]
		item.AppliedDiscount = calculateItemDiscount(item)
		subtotal += item.UnitPrice * float64(item.Quantity)
		itemDiscountTotal += item.AppliedDiscount
	}

	order.Subtotal = round2(subtotal)
	order.ItemDiscountTotal = round2(itemDiscountTotal)

	// Step 2: calculate total order-level discount
	order.AppliedDiscount = calculateOrderDiscount(order)
	order.DiscountTotal = round2(order.ItemDiscountTotal + order.AppliedDiscount)

	// Step 3: allocate order-level discount proportionally
	applyOrderDiscountToItems(order)

	// Step 4: calculate tax and totals for each item
	for i := range order.Items {
		calculateTaxAndLineTotal(&order.Items[i])
	}

	// Step 5: aggregate tax total and final order total
	var taxTotal, totalAmount float64
	for _, item := range order.Items {
		taxTotal += item.TaxAmount
		totalAmount += item.Total
	}

	order.TaxTotal = round2(taxTotal)
	order.Total = round2(totalAmount)
}

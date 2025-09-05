package ordertotals

import (
	"testing"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCalculateOrderTotals_NoDiscounts(t *testing.T) {
	order := &model.Order{
		Items: []model.OrderItem{
			{UnitPrice: 100, Quantity: 1, TaxRate: 0.1},
			{UnitPrice: 50, Quantity: 2, TaxRate: 0.05},
		},
	}

	Calculate(order)

	// Subtotal
	assert.Equal(t, 200.0, order.Subtotal) // 100*1 + 50*2

	// No discounts
	assert.Equal(t, 0.0, order.ItemDiscountTotal)
	assert.Equal(t, 0.0, order.AppliedDiscount)
	assert.Equal(t, 0.0, order.DiscountTotal)

	// Tax
	//expectedTax := round2(100*0.1 + 100*0.05)
	// 10 + 5 = 15
	assert.Equal(t, 15.0, order.TaxTotal)

	// Total
	// expectedTotal := 200.0 + 15.0
	assert.Equal(t, 215.0, order.Total)
}

func TestCalculateOrderTotals_PercentageOrderDiscount(t *testing.T) {
	order := &model.Order{
		Items: []model.OrderItem{
			{UnitPrice: 100, Quantity: 2, TaxRate: 0.1}, // 200
			{UnitPrice: 50, Quantity: 1, TaxRate: 0.05}, // 50
		},
		Discount: model.DiscountInfo{
			Type:   model.DiscountPercentage,
			Amount: 10, // 10% of subtotal
		},
	}

	Calculate(order)

	// Subtotal
	assert.Equal(t, 250.0, order.Subtotal)

	// Order-level discount 10% of 250 = 25
	assert.Equal(t, 25.0, order.AppliedDiscount)
	assert.Equal(t, 25.0, order.DiscountTotal)    // no item discounts
	assert.Equal(t, 0.0, order.ItemDiscountTotal) // no item discounts

	// Proportional applied order discount
	assert.Equal(t, 20.0, order.Items[0].AppliedOrderDiscount) // 200/250 * 25
	assert.Equal(t, 5.0, order.Items[1].AppliedOrderDiscount)  // 50/250 * 25

	// Tax
	expectedTax := round2(order.Items[0].TaxAmount + order.Items[1].TaxAmount)
	assert.Equal(t, expectedTax, order.TaxTotal)

	// Total
	expectedTotal := round2(order.Items[0].Total + order.Items[1].Total)
	assert.Equal(t, expectedTotal, order.Total)
}

func TestCalculateOrderTotals_ZeroItems(t *testing.T) {
	order := &model.Order{}

	Calculate(order)

	assert.Equal(t, 0.0, order.Subtotal)
	assert.Equal(t, 0.0, order.ItemDiscountTotal)
	assert.Equal(t, 0.0, order.AppliedDiscount)
	assert.Equal(t, 0.0, order.DiscountTotal)
	assert.Equal(t, 0.0, order.TaxTotal)
	assert.Equal(t, 0.0, order.Total)
	assert.Empty(t, order.Items)
}

func TestCalculateOrderTotals_ItemAndOrderDiscounts(t *testing.T) {
	order := &model.Order{
		Items: []model.OrderItem{
			{
				UnitPrice: 200,
				Quantity:  1,
				TaxRate:   0.1,
				Discount: model.DiscountInfo{
					Type:   model.DiscountFixed,
					Amount: 20,
				},
			},
			{
				UnitPrice: 100,
				Quantity:  1,
				TaxRate:   0.05,
				Discount: model.DiscountInfo{
					Type:   model.DiscountPercentage,
					Amount: 10,
				},
			},
		},
		Discount: model.DiscountInfo{
			Type:   model.DiscountFixed,
			Amount: 30,
		},
	}

	Calculate(order)

	// Check subtotal
	assert.Equal(t, 300.0, order.Subtotal)

	// Item discounts
	assert.Equal(t, 20.0, order.Items[0].AppliedDiscount) // fixed
	assert.Equal(t, 10.0, order.Items[1].AppliedDiscount) // 10% of 100

	// Order discount
	assert.Equal(t, 30.0, order.AppliedDiscount)

	// AppliedOrderDiscount distribution
	assert.Equal(t, 20.0, order.Items[0].AppliedOrderDiscount) // 200/300 * 30
	assert.Equal(t, 10.0, order.Items[1].AppliedOrderDiscount) // 100/300 * 30

	// Tax
	assert.Equal(t, order.Items[0].TaxAmount+order.Items[1].TaxAmount, order.TaxTotal)

	// Total
	assert.Equal(t, order.Items[0].Total+order.Items[1].Total, order.Total)

	// DiscountTotal
	expectedTotatDiscount := order.AppliedDiscount + order.Items[0].AppliedDiscount + order.Items[1].AppliedDiscount
	assert.Equal(t, expectedTotatDiscount, order.DiscountTotal) // 20 + 10 + 30
}

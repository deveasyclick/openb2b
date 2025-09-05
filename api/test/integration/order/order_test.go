package order_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/test/integration/seed"
	"github.com/deveasyclick/openb2b/test/integration/setup"
	"github.com/stretchr/testify/assert"
)

func TestOrderHandlers(t *testing.T) {
	ts := setup.SetupTestServer()
	defer ts.Close()

	db := setup.SetupTestDB()
	seed.ClearOrders(db) // reset DB for controlled testing
	seed.InsertOrders(db)

	seed.ClearProducts(db)
	product := seed.InsertProducts(db)
	// -------------------- CREATE ORDER --------------------
	t.Run("Create order - success", func(t *testing.T) {
		reqBody := dto.CreateOrderDTO{
			CustomerID: 1,
			Items: []dto.CreateOrderItemDTO{
				{VariantID: product.Variants[0].ID, Quantity: 1},
				{VariantID: product.Variants[1].ID, Quantity: 3},
			},
			Delivery: dto.CreateDeliveryInfoDTO{
				Address: dto.AddressRequired{
					Address: "Street 1",
					City:    "City 1",
					State:   "State 1",
					Country: "Country 1",
					Zip:     "02912",
				},
				TransportFare: 10.0,
			},
			Notes: "Notes",
			Discount: dto.CreateDiscountInfoDTO{
				Type:   model.DiscountFixed,
				Amount: 3,
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orders", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var order response.APIResponse[model.Order]
		err = json.NewDecoder(resp.Body).Decode(&order)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, order.Code)
		assert.Equal(t, order.Message, "success")

		assert.Equal(t, int(order.Data.CustomerID), 1)
		assert.Equal(t, order.Data.Notes, "Notes")
		// Order delivery
		assert.Equal(t, float64(order.Data.Delivery.TransportFare), 10.0)
		assert.Equal(t, order.Data.Delivery.Address.Address, "Street 1")
		assert.Equal(t, order.Data.Delivery.Address.City, "City 1")
		assert.Equal(t, order.Data.Delivery.Address.State, "State 1")
		assert.Equal(t, order.Data.Delivery.Address.Country, "Country 1")
		// Order totals
		assert.Equal(t, int(order.Data.Discount.Amount), 3)
		assert.Equal(t, order.Data.Discount.Type, model.DiscountFixed)
		assert.Equal(t, order.Data.AppliedDiscount, float64(3))
		// discount total = applied discount + sum of all item discounts
		assert.Equal(t, order.Data.DiscountTotal, order.Data.AppliedDiscount+order.Data.Items[0].AppliedDiscount+order.Data.Items[1].AppliedDiscount)
		// item discount total = sum of all item discounts
		assert.Equal(t, order.Data.ItemDiscountTotal, order.Data.Items[0].AppliedDiscount+order.Data.Items[1].AppliedDiscount)
		// order total = sum of all item totals
		assert.Equal(t, order.Data.Total, order.Data.Items[0].Total+order.Data.Items[1].Total)
		// sum of item (unitPrice * qty), before discounts & tax
		subtotal := order.Data.Items[0].UnitPrice*float64(order.Data.Items[0].Quantity) + order.Data.Items[1].UnitPrice*float64(order.Data.Items[1].Quantity)
		assert.Equal(t, order.Data.Subtotal, subtotal)
		// order tax = sum of all item taxes
		assert.Equal(t, order.Data.TaxTotal, order.Data.Items[0].TaxAmount+order.Data.Items[1].TaxAmount)

		// Order items totals
		assert.Equal(t, len(order.Data.Items), 2)
		assert.Equal(t, int(order.Data.Items[0].VariantID), 100)
		assert.Equal(t, order.Data.Items[0].Quantity, 1)
		assert.Equal(t, order.Data.Items[0].UnitPrice, 10.0)
	})

	t.Run("Create order - (zero discount) success", func(t *testing.T) {
		reqBody := dto.CreateOrderDTO{
			CustomerID: 1,
			Items: []dto.CreateOrderItemDTO{
				{VariantID: 99, Quantity: 1},
			},
			Delivery: dto.CreateDeliveryInfoDTO{
				Address: dto.AddressRequired{
					Address: "Street 1",
					City:    "City 1",
					State:   "State 1",
					Country: "Country 1",
					Zip:     "02912",
				},
				TransportFare: 10.0,
			},
			Notes: "Notes",
			Discount: dto.CreateDiscountInfoDTO{
				Type:   model.DiscountPercentage,
				Amount: 0,
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orders", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var order response.APIResponse[model.Order]
		err = json.NewDecoder(resp.Body).Decode(&order)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, order.Code)
		assert.Equal(t, order.Message, "success")

		assert.Equal(t, int(order.Data.CustomerID), 1)
		assert.Equal(t, int(order.Data.Discount.Amount), 0)
		assert.Equal(t, order.Data.Discount.Type, model.DiscountPercentage)
		assert.Equal(t, float64(order.Data.Delivery.TransportFare), 10.0)
		assert.Equal(t, order.Data.Delivery.Address.Address, "Street 1")
		assert.Equal(t, order.Data.Delivery.Address.City, "City 1")
		assert.Equal(t, order.Data.Delivery.Address.State, "State 1")
		assert.Equal(t, order.Data.Delivery.Address.Country, "Country 1")
		assert.Equal(t, order.Data.Notes, "Notes")
		assert.Equal(t, len(order.Data.Items), 1)
		assert.Equal(t, order.Data.Items[0].Quantity, 1)
		assert.Equal(t, order.Data.Items[0].UnitPrice, 20.0)
		assert.Equal(t, int(order.Data.Items[0].VariantID), 99)
		assert.Equal(t, int(order.Data.Total), 20)
	})

	t.Run("Create order - missing name (400)", func(t *testing.T) {
		reqBody := dto.CreateOrderDTO{
			Items: []dto.CreateOrderItemDTO{
				{VariantID: 1, Quantity: 1},
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orders", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- GET PRODUCT --------------------
	t.Run("Get order - success", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orders/1")
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get order - not found (404)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orders/999")
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get order - invalid ID (400)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orders/abc")
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- UPDATE PRODUCT --------------------

	t.Run("Update order - (replace order items) success", func(t *testing.T) {
		note := "This is a note"
		reqBody := dto.UpdateOrderDTO{
			Notes: &note,
			Items: []*dto.CreateOrderItemDTO{
				{Quantity: 3, VariantID: product.Variants[0].ID},
				{Quantity: 3, VariantID: product.Variants[1].ID},
			},
			Discount: &dto.CreateDiscountInfoDTO{
				Type:   model.DiscountFixed,
				Amount: 3,
			},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/orders/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var order response.APIResponse[model.Order]
		err = json.NewDecoder(resp.Body).Decode(&order)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)

		assert.Equal(t, order.Data.Notes, "This is a note")
		// Order totals
		assert.Equal(t, int(order.Data.Discount.Amount), 3)
		assert.Equal(t, order.Data.Discount.Type, model.DiscountFixed)
		assert.Equal(t, order.Data.AppliedDiscount, float64(3))
		// discount total = applied discount + sum of all item discounts
		assert.Equal(t, order.Data.DiscountTotal, order.Data.AppliedDiscount+order.Data.Items[0].AppliedDiscount+order.Data.Items[1].AppliedDiscount)
		// item discount total = sum of all item discounts
		assert.Equal(t, order.Data.ItemDiscountTotal, order.Data.Items[0].AppliedDiscount+order.Data.Items[1].AppliedDiscount)
		// order total = sum of all item totals
		assert.Equal(t, order.Data.Total, order.Data.Items[0].Total+order.Data.Items[1].Total)
		// sum of item (unitPrice * qty), before discounts & tax
		subtotal := order.Data.Items[0].UnitPrice*float64(order.Data.Items[0].Quantity) + order.Data.Items[1].UnitPrice*float64(order.Data.Items[1].Quantity)
		assert.Equal(t, order.Data.Subtotal, subtotal)
		// order tax = sum of all item taxes
		assert.Equal(t, order.Data.TaxTotal, order.Data.Items[0].TaxAmount+order.Data.Items[1].TaxAmount)

		// Order items totals
		assert.Equal(t, len(order.Data.Items), 2)
		assert.Equal(t, int(order.Data.Items[0].VariantID), 100)
		assert.Equal(t, order.Data.Items[0].Quantity, 3)
		assert.Equal(t, order.Data.Items[1].Quantity, 3)
	})

	t.Run("Update order - success", func(t *testing.T) {
		reqBody := map[string]any{"Notes": "This is a note"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/orders/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var order response.APIResponse[model.Order]
		err = json.NewDecoder(resp.Body).Decode(&order)
		if err != nil {
			t.Fatal(err)
		}

		assert.NoError(t, err)
		assert.Equal(t, order.Data.Notes, "This is a note")
		assert.Equal(t, order.Data.OrderNumber, "ORD-123")
		assert.Equal(t, len(order.Data.Items), 2)
	})

	t.Run("Update order - (duplicate order items) success", func(t *testing.T) {
		reqBody := map[string]any{"Notes": "This is a note", "Items": []map[string]any{{"Quantity": 3, "VariantId": product.Variants[1].ID}, {"Quantity": 2, "VariantId": product.Variants[1].ID}}}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/orders/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "response status code")
	})

	t.Run("Update order - invalid ID (400)", func(t *testing.T) {
		reqBody := map[string]any{"name": "Updated Order"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/orders/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update order - not in pending status (400)", func(t *testing.T) {
		reqBody := map[string]any{"Notes": "Updated Note"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/orders/%d", ts.URL, seed.NonPendingOrderId), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var order response.APIResponse[model.Order]
		err = json.NewDecoder(resp.Body).Decode(&order)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, order.Code)
		assert.Equal(t, "cannot update order not in pending: processing", order.Message)
	})

	// -------------------- DELETE PRODUCT --------------------
	t.Run("Delete order - success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orders/1", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete order - not found (404)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orders/999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete order - invalid ID (400)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orders/abc", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

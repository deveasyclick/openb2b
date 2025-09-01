package order_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/modules/order"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/test/integration/setup"
	"github.com/deveasyclick/openb2b/test/integration/setup/seed"
	"github.com/stretchr/testify/assert"
)

func TestOrderHandlers(t *testing.T) {
	ts := setup.SetupTestServer()
	defer ts.Close()

	db := setup.SetupTestDB()
	seed.Clear(db) // reset DB for controlled testing
	Insert(db)
	// -------------------- CREATE ORDER --------------------
	t.Run("Create order - success", func(t *testing.T) {
		reqBody := order.CreateOrderDTO{
			CustomerID: 1,
			Items: []order.CreateOrderItemDTO{
				{VariantID: 1, Quantity: 1, UnitPrice: 10.0, ProductID: 1},
			},
			Delivery: order.CreateDeliveryInfoDTO{
				Address:       &model.Address{Address: "Street 1", City: "City 1", State: "State 1", Country: "Country 1"},
				TransportFare: 10.0,
			},
			Notes: "Notes",
			Discount: order.CreateDiscountInfoDTO{
				Type:   model.DiscountPercentage,
				Amount: 3,
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orders", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var order response.APIResponse[model.Order]
		err = json.NewDecoder(resp.Body).Decode(&order)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, order.Code)
		assert.Equal(t, order.Message, "success")

		assert.Equal(t, int(order.Data.CustomerID), 1)
		assert.Equal(t, int(order.Data.Discount.Amount), 3)
		assert.Equal(t, order.Data.Discount.Type, model.DiscountPercentage)
		assert.Equal(t, float64(order.Data.Delivery.TransportFare), 10.0)
		assert.Equal(t, order.Data.Delivery.Address.Address, "Street 1")
		assert.Equal(t, order.Data.Delivery.Address.City, "City 1")
		assert.Equal(t, order.Data.Delivery.Address.State, "State 1")
		assert.Equal(t, order.Data.Delivery.Address.Country, "Country 1")
		assert.Equal(t, order.Data.Notes, "Notes")
		assert.Equal(t, len(order.Data.Items), 1)
		assert.Equal(t, order.Data.Items[0].Quantity, 1)
		assert.Equal(t, order.Data.Items[0].UnitPrice, 10.0)
		assert.Equal(t, int(order.Data.Items[0].VariantID), 1)
	})

	t.Run("Create order - missing name (400)", func(t *testing.T) {
		reqBody := order.CreateOrderDTO{
			Items: []order.CreateOrderItemDTO{
				{VariantID: 1, Quantity: 1},
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orders", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- GET PRODUCT --------------------
	t.Run("Get order - success", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orders/1")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get order - not found (404)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orders/999")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get order - invalid ID (400)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orders/abc")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- UPDATE PRODUCT --------------------
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
	})

	t.Run("Update order - invalid ID (400)", func(t *testing.T) {
		reqBody := map[string]any{"name": "Updated Order"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/orders/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update order - not in pending status (400)", func(t *testing.T) {
		reqBody := map[string]any{"Notes": "Updated Note"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/orders/%d", ts.URL, nonPendingOrderId), bytes.NewBuffer(body))
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
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete order - not found (404)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orders/999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete order - invalid ID (400)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orders/abc", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

package product_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/test/integration/seed"
	"github.com/deveasyclick/openb2b/test/integration/setup"
	"github.com/stretchr/testify/assert"
)

func TestProductVariantHandlers(t *testing.T) {
	ts := setup.SetupTestServer()
	defer ts.Close()

	db := setup.SetupTestDB()
	seed.InsertProducts(db) // seed initial product(s)

	productID := 1 // assume the first seeded product has ID 1

	t.Run("Create variant success", func(t *testing.T) {
		reqBody := map[string]any{
			"sku":   "SKU-003",
			"color": "Red",
			"size":  "M",
			"price": 19.99,
			"stock": 100,
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+fmt.Sprintf("/api/v1/products/%d/variants", productID), "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Create variant duplicate SKU (409)", func(t *testing.T) {
		reqBody := map[string]any{
			"sku":   "SKU-003",
			"color": "Blue",
			"size":  "L",
			"price": 29.99,
			"stock": 50,
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+fmt.Sprintf("/api/v1/products/%d/variants", productID), "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode) // SKU already exists
	})

	t.Run("Get variant success", func(t *testing.T) {
		resp, err := http.Get(ts.URL + fmt.Sprintf("/api/v1/products/%d/variants/1", productID))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update variant success", func(t *testing.T) {
		reqBody := map[string]any{"price": 24.99}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/1", productID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update variant invalid ID", func(t *testing.T) {
		reqBody := map[string]any{"price": 24.99}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/999", productID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Delete variant success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/1", productID), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete variant not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/999", productID), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

package product_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/internal/shared/dto"
	"github.com/deveasyclick/openb2b/test/integration/setup"
	"github.com/deveasyclick/openb2b/test/integration/setup/seed"
	"github.com/stretchr/testify/assert"
)

func TestProductHandlers(t *testing.T) {
	ts := setup.SetupTestServer()
	defer ts.Close()

	db := setup.SetupTestDB()
	seed.Clear(db) // reset DB for controlled testing

	// -------------------- CREATE PRODUCT --------------------
	t.Run("Create product - success", func(t *testing.T) {
		reqBody := dto.CreateProductDTO{
			Name: "New Product",
			Variants: []dto.CreateProductVariantDTO{
				{SKU: "SKU1", Price: 10.5, Stock: 5},
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/products", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Create product - missing name (400)", func(t *testing.T) {
		reqBody := dto.CreateProductDTO{
			Variants: []dto.CreateProductVariantDTO{
				{SKU: "SKU2", Price: 12, Stock: 3},
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/products", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create product - duplicate name (409)", func(t *testing.T) {
		reqBody := dto.CreateProductDTO{
			Name: "Duplicate Product",
			Variants: []dto.CreateProductVariantDTO{
				{SKU: "SKU3", Price: 15, Stock: 2},
			},
		}
		body, _ := json.Marshal(reqBody)
		resp1, err := http.Post(ts.URL+"/api/v1/products", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp1.StatusCode)

		resp2, err := http.Post(ts.URL+"/api/v1/products", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp2.StatusCode)
	})

	// -------------------- GET PRODUCT --------------------
	t.Run("Get product - success", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/products/1")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get product - not found (404)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/products/999")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get product - invalid ID (400)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/products/abc")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- UPDATE PRODUCT --------------------
	t.Run("Update product - success", func(t *testing.T) {
		reqBody := map[string]any{"name": "Updated Product"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/products/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update product - invalid ID (400)", func(t *testing.T) {
		reqBody := map[string]any{"name": "Updated Product"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/products/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- DELETE PRODUCT --------------------
	t.Run("Delete product - success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/products/1", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete product - not found (404)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/products/999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete product - invalid ID (400)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/products/abc", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

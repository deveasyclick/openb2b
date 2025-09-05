package product_test

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

func TestProductVariantHandlers(t *testing.T) {
	ts := setup.SetupTestServer()
	defer ts.Close()

	db := setup.SetupTestDB()
	product := seed.InsertProducts(db) // seed initial product(s)

	t.Run("Create variant success", func(t *testing.T) {
		reqBody := dto.CreateProductVariantDTO{
			SKU:     "SKU-003",
			Color:   "Red",
			Size:    "M",
			Price:   19.99,
			Stock:   100,
			TaxRate: 0.1,
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+fmt.Sprintf("/api/v1/products/%d/variants", product.ID), "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var variant response.APIResponse[model.Variant]
		err = json.NewDecoder(resp.Body).Decode(&variant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, variant.Code)
		assert.Equal(t, variant.Message, "success")
		assert.Equal(t, variant.Data.SKU, "SKU-003")
		assert.Equal(t, variant.Data.Color, "Red")
		assert.Equal(t, variant.Data.Size, "M")
		assert.Equal(t, variant.Data.Price, 19.99)
		assert.Equal(t, variant.Data.Stock, 100)
		assert.Equal(t, variant.Data.TaxRate, 0.1)
		assert.Equal(t, variant.Data.ProductID, product.ID)

		var updatedProduct model.Product
		err = db.Preload("Variants").First(&updatedProduct, "id = ?", product.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, len(updatedProduct.Variants), 3)
	})

	t.Run("Create variant duplicate SKU (409)", func(t *testing.T) {
		reqBody := dto.CreateProductVariantDTO{
			SKU:     "SKU-003",
			Color:   "Blue",
			Size:    "L",
			Price:   29.99,
			Stock:   50,
			TaxRate: 0.1,
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+fmt.Sprintf("/api/v1/products/%d/variants", product.ID), "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode) // SKU already exists
	})

	t.Run("Get variant success", func(t *testing.T) {
		resp, err := http.Get(ts.URL + fmt.Sprintf("/api/v1/products/%d/variants/101", product.ID))
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update variant success", func(t *testing.T) {
		reqBody := map[string]any{"price": 24.99}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/99", product.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var variant response.APIResponse[model.Variant]
		err = json.NewDecoder(resp.Body).Decode(&variant)
		assert.NoError(t, err)
		assert.Equal(t, variant.Data.Price, 24.99)
		assert.Equal(t, variant.Data.ProductID, product.ID)
	})

	t.Run("Update variant invalid ID", func(t *testing.T) {
		reqBody := map[string]any{"price": 24.99}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/999", product.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete variant success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/101", product.ID), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete variant not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+fmt.Sprintf("/api/v1/products/%d/variants/999", product.ID), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

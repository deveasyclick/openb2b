package customer_test

import (
	"bytes"
	"encoding/json"
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
	seed.ClearCustomers(db) // reset DB for controlled testing

	// -------------------- CREATE CUSTOMER --------------------
	t.Run("Create customer - success", func(t *testing.T) {
		reqBody := dto.CreateCustomerDTO{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "i6oP7@example.com",
			Address: &dto.AddressOptional{
				Address: "123 Market Street",
				City:    "San Francisco",
				State:   "California",
				Country: "USA",
				Zip:     "02912",
			},
			Company:     "OpenB2B",
			PhoneNumber: "+1-202-555-0199",
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/customers", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var customer response.APIResponse[model.Customer]
		err = json.NewDecoder(resp.Body).Decode(&customer)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, customer.Code)
		assert.Equal(t, customer.Message, "success")

		assert.Equal(t, int(customer.Data.ID), 1)
		assert.Equal(t, customer.Data.FirstName, "John")
		assert.Equal(t, customer.Data.LastName, "Doe")
		assert.Equal(t, customer.Data.Email, "i6oP7@example.com")
		assert.Equal(t, customer.Data.PhoneNumber, "+1-202-555-0199")
		assert.Equal(t, customer.Data.Address.Address, "123 Market Street")
		assert.Equal(t, customer.Data.Address.City, "San Francisco")
		assert.Equal(t, customer.Data.Address.State, "California")
		assert.Equal(t, customer.Data.Address.Country, "USA")
		assert.Equal(t, customer.Data.Company, "OpenB2B")
	})

	t.Run("Create customer - (no optional fields) success", func(t *testing.T) {
		reqBody := dto.CreateCustomerDTO{
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+1-202-555-0199",
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/customers", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var customer response.APIResponse[model.Customer]
		err = json.NewDecoder(resp.Body).Decode(&customer)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, customer.Code)
		assert.Equal(t, customer.Message, "success")

		assert.Equal(t, int(customer.Data.ID), 2)
		assert.Equal(t, customer.Data.FirstName, "John")
		assert.Equal(t, customer.Data.LastName, "Doe")
		assert.Equal(t, customer.Data.Email, "")
		assert.Equal(t, customer.Data.PhoneNumber, "+1-202-555-0199")
		assert.Equal(t, customer.Data.Company, "")
		assert.Equal(t, customer.Data.Address, (*model.Address)(nil))
	})

	t.Run("Create customer - missing required fields (400)", func(t *testing.T) {
		reqBody := dto.CreateCustomerDTO{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "i6oP7@example.com",
			Company:   "OpenB2B",
			Address: &dto.AddressOptional{
				Address: "123 Market Street",
				City:    "San Francisco",
				State:   "California",
				Country: "USA",
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/customers", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- GET PRODUCT --------------------
	t.Run("Get customer - success", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/customers/1")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get customer - not found (404)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/customers/999")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get customer - invalid ID (400)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/customers/abc")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- UPDATE PRODUCT --------------------

	t.Run("Update customer - success", func(t *testing.T) {
		reqBody := dto.CreateCustomerDTO{
			FirstName:   "Johny",
			LastName:    "Darwin",
			Email:       "test@example.com",
			Company:     "OpenC2C",
			PhoneNumber: "+234623370870",
			Address: &dto.AddressOptional{
				Address: "washington avenue",
				City:    "Houston",
				State:   "Texas",
				Country: "United Kingdom",
				Zip:     "123456",
			},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/customers/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var customer response.APIResponse[model.Customer]
		err = json.NewDecoder(resp.Body).Decode(&customer)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)

		assert.Equal(t, int(customer.Data.ID), 1)
		assert.Equal(t, customer.Data.FirstName, "Johny")
		assert.Equal(t, customer.Data.LastName, "Darwin")
		assert.Equal(t, customer.Data.Email, "test@example.com")
		assert.NotEqual(t, customer.Data.PhoneNumber, "+234623370870") // should not update phone number
		assert.Equal(t, customer.Data.Company, "OpenC2C")
		assert.Equal(t, customer.Data.Address.Address, "washington avenue")
		assert.Equal(t, customer.Data.Address.City, "Houston")
		assert.Equal(t, customer.Data.Address.State, "Texas")
		assert.Equal(t, customer.Data.Address.Country, "United Kingdom")
		assert.Equal(t, customer.Data.Address.Zip, "123456")
	})

	t.Run("Update customer - invalid ID (400)", func(t *testing.T) {
		reqBody := map[string]any{"name": "Updated Customer"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/customers/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- DELETE PRODUCT --------------------
	t.Run("Delete customer - success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/customers/1", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete customer - not found (404)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/customers/999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete customer - invalid ID (400)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/customers/abc", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

package org_test

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
	seed.ClearOrgs(db) // reset DB for controlled testing

	// -------------------- CREATE ORG --------------------
	t.Run("Create org- success", func(t *testing.T) {
		reqBody := dto.CreateOrgDTO{
			Name:             "OpenB2B",
			Logo:             "https://www.openb2b.com/logo.png",
			OrganizationName: "OpenB2B NG",
			OrganizationUrl:  "https://www.openb2b.com",
			Email:            "info@openb2b.com",
			Phone:            "+1-202-555-0199",
			Address: dto.AddressRequired{
				Address: "123 Market Street",
				City:    "San Francisco",
				State:   "California",
				Country: "USA",
				Zip:     "02912",
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orgs", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var org response.APIResponse[model.Org]
		err = json.NewDecoder(resp.Body).Decode(&org)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, org.Code)
		assert.Equal(t, org.Message, "success")

		assert.Equal(t, int(org.Data.ID), 1)
		assert.Equal(t, org.Data.OrganizationName, "OpenB2B NG")
		assert.Equal(t, org.Data.OrganizationUrl, "https://www.openb2b.com")
		assert.Equal(t, org.Data.Logo, "https://www.openb2b.com/logo.png")
		assert.Equal(t, org.Data.Name, "OpenB2B")
		assert.Equal(t, org.Data.Email, "info@openb2b.com")
		assert.Equal(t, org.Data.Phone, "+1-202-555-0199")
		assert.Equal(t, org.Data.Address.Address, "123 Market Street")
		assert.Equal(t, org.Data.Address.City, "San Francisco")
		assert.Equal(t, org.Data.Address.State, "California")
		assert.Equal(t, org.Data.Address.Country, "USA")
	})

	t.Run("Create org - (no optional fields) success", func(t *testing.T) {
		reqBody := dto.CreateOrgDTO{
			Name:             "OpenB2C",
			OrganizationName: "OpenB2B NG",
			Phone:            "+1-202-555-0199",
			Email:            "info@openb2c.com",
			Address: dto.AddressRequired{
				Address: "123 Market Street",
				City:    "San Francisco",
				State:   "California",
				Country: "USA",
				Zip:     "02912",
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orgs", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var org response.APIResponse[model.Org]
		err = json.NewDecoder(resp.Body).Decode(&org)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, org.Code)
		assert.Equal(t, org.Message, "success")

		assert.Equal(t, int(org.Data.ID), 2)
		assert.Equal(t, org.Data.OrganizationName, "OpenB2B NG")
		assert.Equal(t, org.Data.Name, "OpenB2C")
		assert.Equal(t, org.Data.Email, "info@openb2c.com")
		assert.Equal(t, org.Data.Phone, "+1-202-555-0199")
		assert.Equal(t, org.Data.Address.Address, "123 Market Street")
		assert.Equal(t, org.Data.Address.City, "San Francisco")
		assert.Equal(t, org.Data.Address.State, "California")
		assert.Equal(t, org.Data.Address.Country, "USA")
	})

	t.Run("Create org - missing required fields (400)", func(t *testing.T) {
		reqBody := dto.CreateOrgDTO{}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/orgs", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- GET PRODUCT --------------------
	t.Run("Get org - success", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orgs/1")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get org - not found (404)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orgs/999")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get org - invalid ID (400)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/orgs/abc")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- UPDATE PRODUCT --------------------

	t.Run("Update org - success", func(t *testing.T) {
		reqBody := dto.UpdateOrgDTO{
			Name:             "OpenC2C",
			OrganizationName: "OpenC2C NG", // should not update
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/orgs/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var org response.APIResponse[model.Org]
		err = json.NewDecoder(resp.Body).Decode(&org)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)

		assert.Equal(t, int(org.Data.ID), 1)
		assert.Equal(t, org.Data.Name, "OpenC2C")
		assert.Equal(t, org.Data.Email, "info@openb2b.com")
		assert.Equal(t, org.Data.OrganizationName, "OpenB2B NG")
		assert.Equal(t, org.Data.Address.Address, "123 Market Street")
		assert.Equal(t, org.Data.Address.City, "San Francisco")
		assert.Equal(t, org.Data.Address.State, "California")
		assert.Equal(t, org.Data.Address.Country, "USA")
		assert.Equal(t, org.Data.Address.Zip, "02912")
	})

	t.Run("Update org - invalid ID (400)", func(t *testing.T) {
		reqBody := map[string]any{"name": "Updated Org"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/orgs/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// -------------------- DELETE PRODUCT --------------------
	t.Run("Delete org - success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orgs/1", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete org - not found (404)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orgs/209", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete org - invalid ID (400)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orgs/abc", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

package customer_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/response"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/test/integration/seed"
	"github.com/deveasyclick/openb2b/test/integration/setup"
	"github.com/stretchr/testify/assert"
)

const email = "i6oP7@example.com"

func TestOrderHandlers(t *testing.T) {
	ts := setup.SetupTestServer()
	defer ts.Close()

	db := setup.SetupTestDB()
	seed.ClearUsers(db) // reset DB for controlled testing

	// -------------------- CREATE CUSTOMER --------------------
	t.Run("Create user - success", func(t *testing.T) {
		data := map[string]interface{}{
			"id":              "user_1",
			"first_name":      "John",
			"last_name":       "Doe",
			"email_addresses": []types.ClerkEmail{{ID: "email_1", EmailAddress: email}},
		}

		reqBody := types.WebhookEvent{
			Type: "user.created",
			Data: data,
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/webhooks/handleEvents", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var user model.User
		err = db.First(&user, "email = ?", email).Error
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, user.FirstName, "John")
		assert.Equal(t, user.LastName, "Doe")
	})

	t.Run("Create user - error (missing email)", func(t *testing.T) {

		data := map[string]interface{}{
			"id":              "user_1",
			"first_name":      "John",
			"last_name":       "Doe",
			"email_addresses": []types.ClerkEmail{},
		}

		reqBody := types.WebhookEvent{
			Type: "user.created",
			Data: data,
		}
		body, _ := json.Marshal(reqBody)
		resp, err := http.Post(ts.URL+"/api/v1/webhooks/handleEvents", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res response.APIResponse[string]
		err = json.NewDecoder(resp.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, res.Data, apperrors.ErrEmailNotFoundInClerkWebhook)
	})

	t.Run("Create user - success (invalid body)", func(t *testing.T) {
		body, _ := json.Marshal(nil)
		resp, err := http.Post(ts.URL+"/api/v1/webhooks/handleEvents", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

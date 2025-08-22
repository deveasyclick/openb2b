package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/deveasyclick/openb2b/pkg/apperrors"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateRequest(r *http.Request, req interface{}) []apperrors.ValidationError {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return []apperrors.ValidationError{{
			Field: "body",
			Tag:   "json",
			Value: err.Error(),
		}}
	}

	var validationErrors []apperrors.ValidationError

	err := validate.Struct(req)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			for _, err := range validationErr {
				validationErrors = append(validationErrors, apperrors.ValidationError{
					Field: err.Field(),
					Tag:   err.Tag(),
					Value: fmt.Sprintf("%v", err.Value()),
				})
			}
		}
	}

	return validationErrors
}

func WriteValidationResponse(w http.ResponseWriter, errors []apperrors.ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": errors,
	})
}

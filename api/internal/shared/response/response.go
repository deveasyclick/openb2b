// Package response provides utilities for writing standardized HTTP responses,
// including JSON-encoded success and error messages. It is intended to simplify
// API response handling across the application and ensure consistency.
package response

import (
	"encoding/json"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

// WriteJSONError writes a JSON-encoded error response to the given http.ResponseWriter
// using the provided *apperrors.APIError. It sets the appropriate HTTP status code
// and encodes the error in a standardized response format.
//
// If the error represents a server-side failure (5xx), the function will also
// log the internal error details using logger. If the response cannot be encoded,
// an additional log entry is written.

func WriteJSONError(w http.ResponseWriter, err *apperrors.APIError, logger interfaces.Logger) {
	if err.IsServerError() {
		logger.Error(err.Message, "err", err.InternalMsg)
	} else {
		logger.Debug("client error", "msg", err.Message, "code", err.Code)
	}

	w.WriteHeader(err.Code)
	if encodeErr := json.NewEncoder(w).Encode(err.ToResponse()); encodeErr != nil {
		logger.Error("failed to encode API error response", "err", encodeErr)
	}
}

// APIResponse is a standard success response wrapper
type APIResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data,omitempty"`
}

// WriteJSONSuccess writes a structured success response
func WriteJSONSuccess(w http.ResponseWriter, statusCode int, data any, logger interfaces.Logger) {
	resp := APIResponse{
		Code:    statusCode,
		Message: "success",
		Data:    data,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error(apperrors.ErrEncodeResponse, "error", err)
	}
}

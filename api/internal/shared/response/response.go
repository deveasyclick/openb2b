// Package response provides utilities for writing standardized HTTP responses,
// including JSON-encoded success and error messages. It is intended to simplify
// API response handling across the application and ensure consistency.
package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deveasyclick/openb2b/pkg/apperrors"
)

// WriteJSONError writes a JSON-encoded error response to the given http.ResponseWriter
// using the provided *apperrors.APIError. It sets the appropriate HTTP status code
// and encodes the error in a standardized response format.
//
// If the error represents a server-side failure (5xx), the function will also
// log the internal error details using slog. If the response cannot be encoded,
// an additional log entry is written.

func WriteJSONError(w http.ResponseWriter, err *apperrors.APIError) {
	if err.IsServerError() {
		slog.Error(err.Message, "err", err.InternalMsg)
	}

	w.WriteHeader(err.Code)
	if encodeErr := json.NewEncoder(w).Encode(err.ToResponse()); encodeErr != nil {
		slog.Error("failed to encode API error response", "err", encodeErr)
	}
}

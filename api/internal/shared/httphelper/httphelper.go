package httphelper

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deveasyclick/openb2b/pkg/apperrors"
)

func WriteJSONError(w http.ResponseWriter, err *apperrors.APIError) {
	if err.IsServerError() {
		slog.Error(err.Message, "err", err.InternalMsg)
	}

	w.WriteHeader(err.Code)
	if encodeErr := json.NewEncoder(w).Encode(err.ToResponse()); encodeErr != nil {
		slog.Error("failed to encode API error response", "err", encodeErr)
	}
}

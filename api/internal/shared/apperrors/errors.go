package apperrors

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// APIErrorResponse represents the JSON error returned by the API
// @Description Error response object
type APIErrorResponse struct {
	// Application-specific error code
	Code int `json:"code" example:"400"`

	// Human-readable message for clients
	Message string `json:"message" example:"invalid request body"`
}

type APIError struct {
	Code        int
	Message     string
	InternalMsg string // detailed message for logs if code is 500
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) IsServerError() bool {
	return e.Code >= 500
}

func (e *APIError) ToResponse() APIErrorResponse {
	return APIErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}

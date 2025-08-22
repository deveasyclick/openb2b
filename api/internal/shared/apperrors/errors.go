package apperrors

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
type APIErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

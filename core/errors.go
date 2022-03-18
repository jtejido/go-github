package core

// ErrorResponse is the standard error replied by the API.
// Note that each route/route group may implement its own replies.
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (er ErrorResponse) Error() string {
	return er.Message
}

// NewErrorResponse returns an ErrorResponse with the specified message.
func NewErrorResponse(str string) *ErrorResponse {
	return &ErrorResponse{Message: str}
}

// NewErrorResponseWithCode returns an ErrorResponse with the specified message and code.
func NewErrorResponseWithCode(str string, code int) *ErrorResponse {
	return &ErrorResponse{Message: str, Code: code}
}

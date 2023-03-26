package resp

// SuccessResponse wraps return data with a message
// @Description Wrap successful API responses with a message and data.
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// ErrorResponse wraps errors with a message
// @Description Wrap failed API responses with a message and error.
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type ErrorBag map[string]string

// FormErrorResponse wraps form errors with a message
// @Description Wrapped failed API responses with an ErrorBag and message.
// @Description Used to create synchronicity with the front end.
type FormErrorResponse struct {
	Errors  ErrorBag `json:"errors"`
	Message string   `json:"message"`
}

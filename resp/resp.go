package resp

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// Sender is a struct that holds the fiber.Ctx needed to send curated API responses.
type Sender struct{ ctx *fiber.Ctx }

// New creates a new Sender struct that can be used to easily send ErrorResponse, FormErrorResponse, and SuccessResponse to a fiber.Ctx
func New(c *fiber.Ctx) Sender {
	return Sender{c}
}

// Success sends a message and data interface with a http.StatusOK status code
func (r Sender) Success(msg string, data interface{}) error {
	_ = r.ctx.SendStatus(http.StatusOK)

	// prevent `null` values in data when arrays are empty
	if fmt.Sprintf("%t", data) == "[]" {
		data = make([]string, 0)
	}
	if data == nil {
		data = fiber.Map{}
	}

	return r.ctx.JSON(SuccessResponse{
		Message: msg,
		Data:    data,
	})
}

// Error sends a message and error message with a http.StatusUnprocessableEntity status code
func (r Sender) Error(msg string, err error) error {
	return r.sendErrorWithStatusCode(msg, err, http.StatusUnprocessableEntity)
}

// FormError sends a message and ErrorBag with a http.StatusUnprocessableEntity status code. Use for UI feedback.
func (r Sender) FormError(msg string, errors ErrorBag) error {
	return r.sendFormErrorWithStatusCode(msg, errors, http.StatusUnprocessableEntity)
}

// Unauthorized is a helper function to send an error with a http.StatusUnauthorized status code
func (r Sender) Unauthorized(msg string, err error) error {
	return r.sendErrorWithStatusCode(msg, err, http.StatusUnauthorized)
}

// BadRequest is a helper function to send an error with a http.StatusBadRequest status code
func (r Sender) BadRequest(msg string, err error) error {
	return r.sendErrorWithStatusCode(msg, err, http.StatusBadRequest)
}

// NotFound is a helper function to send an error with a http.StatusNotFound status code
func (r Sender) NotFound(msg string, err error) error {
	return r.sendErrorWithStatusCode(msg, err, http.StatusNotFound)
}

// sendErrorWithStatusCode sends a message and error with a given status code
func (r Sender) sendErrorWithStatusCode(msg string, err error, code int) error {
	_ = r.ctx.SendStatus(code)
	return r.ctx.JSON(ErrorResponse{
		Message: msg,
		Error:   err.Error(),
	})
}

// sendFormErrorWithStatusCode sends a message and ErrorBag with a given status code
func (r Sender) sendFormErrorWithStatusCode(msg string, errors ErrorBag, code int) error {
	_ = r.ctx.SendStatus(code)

	return r.ctx.JSON(FormErrorResponse{
		Message: msg,
		Errors:  errors,
	})
}

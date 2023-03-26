package trex

import (
	"github.com/gofiber/fiber/v2"
	"testing"
)

type Map map[string]interface{}

type ITestSuite interface {
	SetupSuite()
}

type ITestable interface {
	T() *testing.T
}

type IBenchmarkable interface {
	B() *testing.B
}

type IFiberTestSuite interface {
	ITestSuite
	ITestable
	App() *fiber.App
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ErrorMap map[string]string
type FormErrorResponse struct {
	Errors ErrorMap `json:"errors"`
}

// Fields get all fields from the FormErrorResponse
func (j FormErrorResponse) Fields() []string {
	var res []string
	for k, _ := range j.Errors {
		res = append(res, k)
	}

	return res
}

// Contains exhaustively checks all fields in FormErrorResponse for matches.
func (j FormErrorResponse) Contains(fields ...string) bool {
	needs := len(fields)
	found := 0
	for _, field := range fields {
		if _, ok := j.Errors[field]; ok {
			found++
		}
	}

	return needs == found
}

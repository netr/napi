package napi

import (
	"github.com/gofiber/fiber/v2"
	"testing"
)

func Test_NewResourceController(t *testing.T) {
	nc := NewResourceController("test")
	if nc.prefix != "test" {
		t.Fatalf("NewResourceController should not return nil")
	}
}

type MockResourceController struct {
	ResourceController
}

func NewMockResourceController() IResourceController {
	return &MockResourceController{
		ResourceController: NewResourceController("test"),
	}
}

func (m *MockResourceController) Index() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("index")
	}
}

func (m *MockResourceController) Show() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("show")
	}
}

func (m *MockResourceController) Store() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("store")
	}
}

func (m *MockResourceController) Update() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("update")
	}
}

func (m *MockResourceController) Destroy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("destroy")
	}
}

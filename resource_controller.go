package napi

import "github.com/gofiber/fiber/v2"

type IResourceController interface {
	Index() fiber.Handler
	Show() fiber.Handler
	Store() fiber.Handler
	Update() fiber.Handler
	Destroy() fiber.Handler
	Prefix() string
}

type ResourceController struct {
	IResourceController
	prefix string
}

func (c *ResourceController) Prefix() string {
	return c.prefix
}

func NewResourceController(prefix string) ResourceController {
	return ResourceController{
		prefix: prefix,
	}
}

package sweets

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type FiberSuite struct {
	app    *fiber.App
	router fiber.Router
}

// NewFiberSuite is used to instantiate a new fiber.App test suite. Typically called in SetupSuite() while testing controllers.
func (suite *FiberSuite) NewFiberSuite(baseUrl string) *fiber.App {
	app := fiber.New()
	api := app.Group(baseUrl)

	suite.app = app
	suite.router = api
	return app
}

// Route is a helper to auto generate URLs based on route names and arguments.
//
// USAGE: suite.Route("route.name", param1, param2) = /category/1/subcategory/2
func (suite *FiberSuite) Route(name string, args ...interface{}) string {
	for _, routes := range suite.App().Stack() {
		for _, route := range routes {
			if route.Name == name {
				path := route.Path
				for i, param := range route.Params {
					path = strings.Replace(path, ":"+param, fmt.Sprintf("%v", args[i]), 1)
				}
				return path
			}
		}
	}

	return ""
}

// Router is a helper to get the underlying fiber.Router
func (suite *FiberSuite) Router() fiber.Router {
	return suite.router
}

// App is a helper to get the underlying *fiber.App
func (suite *FiberSuite) App() *fiber.App {
	return suite.app
}

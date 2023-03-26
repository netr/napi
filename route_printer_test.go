package napi

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func Test_NewRoutePrinter_EmptyFiberApp(t *testing.T) {
	rp := NewRoutePrinter(&fiber.App{})
	require.NotNil(t, rp, "should have got printer struct with empty fiber app")
}

func Test_parseRouteHandler_Success(t *testing.T) {
	val := `hulk/controllers.AuthController.Login.func1`
	ctrl, fn := parseFiberHandler(val)

	require.Equalf(t, "AuthController", ctrl, "wanted AuthController, got: %s", ctrl)
	require.Equalf(t, "Login", fn, "wanted Login, got: %s", fn)
}

func Test_parseRouteHandler_Success_OtherPrefixes(t *testing.T) {
	val := `hulk/server.AuthBaby.GetAllOfTheMoney.func1`
	ctrl, fn := parseFiberHandler(val)

	require.Equalf(t, "AuthBaby", ctrl, "wanted AuthBaby, got: %s", ctrl)
	require.Equalf(t, "GetAllOfTheMoney", fn, "wanted GetAllOfTheMoney, got: %s", fn)
}

func Test_parseRouteHandler_WillParseUglyText(t *testing.T) {
	val := `github.com/gofiber/jwt/v3.New.func1`
	ctrl, fn := parseFiberHandler(val)

	require.Equalf(t, "", ctrl, "wanted '', got: %s", ctrl)
	require.Equalf(t, "", fn, "wanted '', got: %s", fn)
}

func Test_Hydrate_ShouldHaveParams(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api")
	api.Get("/:user/:id", func(ctx *fiber.Ctx) error { return nil })

	bag := NewRoutePrinter(app).Hydrate()

	require.Len(t, bag.items[0].params, 2, "should have two params")
}

func Test_Hydrate_ShouldHaveCorrectPath(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api")
	api.Get("/:user/:id", func(ctx *fiber.Ctx) error { return nil })

	bag := NewRoutePrinter(app).Hydrate()

	require.Equal(t, "/api/:user/:id", bag.items[0].path, "should have correct path")
}

type MockController struct{}

func (r MockController) MockHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func ExamplePrintRoutes() {
	app := fiber.New()
	mc := &MockController{}
	app.Get("/:user/:id", mc.MockHandler())
	app.Get("/test", func(ctx *fiber.Ctx) error { return nil })

	NewRoutePrinter(app).PrintPretty()
}

package napi

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"testing"
)

func TestNewRouter(t *testing.T) {
	f := fiber.New()
	r := NewRouter(f, "/api")

	// check if the prefix router is set up properly
	if r.prefixRouter == nil {
		t.Fatalf("prefix router should not be nil")
	}
}

func TestRouter_Resource(t *testing.T) {
	f := fiber.New()
	r := NewRouter(f, "/api")
	r.Resource("/test", NewMockResourceController())
	expectedNames := []string{"test.index", "test.show", "test.create", "test.update", "test.destroy"}

	if r.app.Stack() == nil {
		t.Fatalf("app stack should not be nil")
	}

	for _, routes := range r.app.Stack() {
		for _, route := range routes {
			if route.Method != "HEAD" {
				if !strings.HasPrefix(route.Path, "/api/test") {
					t.Fatalf("route path should start with /api/test, got: %s", route.Path)
				}

				found := false
				for _, name := range expectedNames {
					if route.Name == name {
						found = true
						expectedNames = append(expectedNames[:0], expectedNames[1:]...)
						break
					}
				}

				if !found {
					t.Fatalf("route name should be one of %v, got: %s", expectedNames, route.Name)
				}
			}
		}
	}

	if len(expectedNames) != 0 {
		t.Fatalf("expected no more names, got: %v", expectedNames)
	}
}

func TestRouter_AddGroup(t *testing.T) {
	expectedNames := []string{"boom.index", "boom.test"}
	rg := NewRouteGroup("/test", "boom")
	rg.AddRoutes(
		NewRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error { return nil }, "index"),
		NewRoute(fiber.MethodGet, "/test", func(c *fiber.Ctx) error { return nil }, "test"),
	)

	f := fiber.New()
	r := NewRouter(f, "/api")
	r.AddGroup(rg)

	if r.app.Stack() == nil {
		t.Fatalf("app stack should not be nil")
	}

	for _, routes := range r.app.Stack() {
		for _, route := range routes {
			if route.Method != "HEAD" {
				if !strings.HasPrefix(route.Path, "/api/test") {
					t.Fatalf("route path should start with /api/test, got: %s", route.Path)
				}

				found := false
				for _, name := range expectedNames {
					if route.Name == name {
						found = true
						expectedNames = append(expectedNames[:0], expectedNames[1:]...)
						break
					}
				}

				if !found {
					t.Fatalf("route name should be one of %v, got: %s", expectedNames, route.Name)
				}
			}
		}
	}

	if len(expectedNames) != 0 {
		t.Fatalf("expected no more names, got: %v", expectedNames)
	}
}

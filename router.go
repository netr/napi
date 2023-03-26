package napi

import "github.com/gofiber/fiber/v2"

// Router holds the underlying fiber.App used to set up your API routes.
type Router struct {
	app          *fiber.App
	prefixRouter fiber.Router
}

// NewRouter instantiates a new Routes struct used to set up your API routes.
// urlPrefix is the prefix used for all routes.
func NewRouter(app *fiber.App, urlPrefix string) Router {
	if urlPrefix == "" {
		urlPrefix = "/"
	}
	if len(urlPrefix) > 1 && urlPrefix[len(urlPrefix)-1] != '/' {
		urlPrefix += "/"
	}

	return Router{
		app:          app,
		prefixRouter: app.Group(urlPrefix),
	}
}

// AddGroup adds a new group of routes to the router
func (r *Router) AddGroup(group *RouteGroup) {
	g := r.prefixRouter.Group(group.urlPrefix)

	for _, route := range group.routes {
		name := route.Name
		if group.namePrefix != "" {
			name = group.namePrefix + "." + name
		}
		g.Add(route.Method, route.Path, route.Handler).Name(name)
	}
}

// Resource creates a new resource controller with the given path prefix and optional
func (r *Router) Resource(urlPrefix string, controller IResourceController) *Router {
	g := r.prefixRouter.Group(urlPrefix)

	g.Get("/", controller.Index()).Name(controller.Prefix() + ".index")
	g.Get("/:id", controller.Show()).Name(controller.Prefix() + ".show")
	g.Post("/", controller.Store()).Name(controller.Prefix() + ".store")
	g.Patch("/:id", controller.Update()).Name(controller.Prefix() + ".update")
	g.Delete("/:id", controller.Destroy()).Name(controller.Prefix() + ".destroy")
	return r
}

// RouteGroup is a group of routes
type RouteGroup struct {
	routes     []*GroupRoute
	namePrefix string
	urlPrefix  string
}

// NewRouteGroup creates a new route group
func NewRouteGroup(urlPrefix, namePrefix string) *RouteGroup {
	return &RouteGroup{
		urlPrefix:  urlPrefix,
		namePrefix: namePrefix,
	}
}

// AddRoute adds a new route to the router
func (r *RouteGroup) AddRoute(route *GroupRoute) {
	r.routes = append(r.routes, route)
}

// AddRoutes adds a new group of routes to the router
func (r *RouteGroup) AddRoutes(routes ...*GroupRoute) {
	r.routes = append(r.routes, routes...)
}

// GroupRoute is a route in a group
type GroupRoute struct {
	Method  string
	Path    string
	Handler fiber.Handler
	Name    string
}

// NewRoute creates a new route in a group
func NewRoute(method, path string, handler fiber.Handler, name string) *GroupRoute {
	return &GroupRoute{
		Method:  method,
		Path:    path,
		Handler: handler,
		Name:    name,
	}
}

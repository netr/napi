package ctrl

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Routes holds the underlying fiber.App used to set up your API routes.
type Routes struct{ app *fiber.App }

// NewRoutes instantiates a new Routes struct used to set up your API routes.
func NewRoutes(app *fiber.App) Routes {
	return Routes{app: app}
}

// Setup initializes all the controllers and their routes
func (r Routes) Setup(db *gorm.DB) {
	r.Accounts(NewAccountController(db))
}

// Accounts initializes the AccountController routes
func (r Routes) Accounts(c *AccountController) {
	api := r.app.Group("/accounts")
	api.Get("/", c.Index()).Name("accounts.index")
	api.Post("/", c.Store()).Name("accounts.store")
	api.Post("/:id", c.Update()).Name("accounts.update")
}

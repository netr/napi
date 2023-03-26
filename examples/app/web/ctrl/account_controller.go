package ctrl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/netr/napi"
	"github.com/netr/napi/examples/app/db/models"
	"github.com/netr/napi/examples/app/db/repos"
	"github.com/netr/napi/examples/app/web/dto"
	"github.com/netr/napi/resp"
	"gorm.io/gorm"
)

type AccountController struct {
	db   *gorm.DB
	repo *repos.AccountRepo
	napi.CanValidate
}

func NewAccountController(db *gorm.DB) *AccountController {
	return &AccountController{
		db:   db,
		repo: repos.NewAccountRepo(db),
	}
}

func (ac *AccountController) Index() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accs, err := ac.repo.GetAll()
		if err != nil {
			return resp.New(c).Error("get all accounts", err)
		}

		return resp.New(c).Success("success", accs)
	}
}

func (ac *AccountController) Store() fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(dto.AccountStoreRequest)
		if err := ac.Validate(c, request); err != nil {
			return resp.New(c).FormError("creating account", err.Error())
		}

		var err error
		model := new(models.Account)
		if model, err = ac.repo.Create(&models.Account{
			Username: request.Username,
			Password: request.Password,
		}); err != nil {
			return err
		}

		return resp.New(c).Success("success", model)
	}
}

func (ac *AccountController) Update() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(dto.AccountUpdateRequest)
		if err := ac.Validate(c, request); err != nil {
			return resp.New(c).FormError("updating account", err.Error())
		}

		var err error
		model := new(models.Account)
		if model, err = ac.repo.Update(id, request.Password); err != nil {
			return err
		}

		return resp.New(c).Success("success", model)
	}
}

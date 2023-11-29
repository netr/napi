package repos

import (
	"github.com/netr/napi"
	"github.com/netr/napi/examples/app/db/models"
	"gorm.io/gorm"
)

type AccountRepo struct {
	db napi.IRepository[*gorm.DB]
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{
		db: napi.NewGormRepository(db),
	}
}

func (a *AccountRepo) GetAll() ([]models.Account, error) {
	var accounts []models.Account
	if tx := a.db.DB().Model(&models.Account{}).Find(&accounts); tx.Error != nil {
		return nil, tx.Error
	}

	return accounts, nil
}

func (a *AccountRepo) Create(acc *models.Account) (*models.Account, error) {
	if err := a.db.Create(acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (a *AccountRepo) Update(id interface{}, password string) (*models.Account, error) {
	model := new(models.Account)
	if err := a.db.Update(model, id, napi.UpdateMap{"password": password}); err != nil {
		return nil, err
	}

	return model, nil
}

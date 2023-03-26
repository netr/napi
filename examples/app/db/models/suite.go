package models

import (
	"github.com/netr/napi/factory"
	"gorm.io/gorm"
)

type FactorySuite struct {
	factory *factory.Factory
}

func (suite *FactorySuite) Factory() *factory.Factory {
	return suite.factory
}

func (suite *FactorySuite) NewFactorySuite(db *gorm.DB) {
	suite.factory = factory.New(db).Add(&Account{})
}

func (suite *FactorySuite) CreateAccount() *Account {
	return suite.Factory().Create(Account{}).(*Account)
}

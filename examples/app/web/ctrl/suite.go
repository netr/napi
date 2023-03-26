package ctrl

import (
	"github.com/netr/napi/examples/app/db/models"
	"github.com/netr/napi/sweets"
)

type ControllerSuite struct {
	sweets.Suite
	sweets.FiberSuite
	sweets.GormSuite
	models.FactorySuite
}

func (suite *ControllerSuite) SetupSuite() {
	suite.NewFiberSuite("/")
	suite.NewGormSuite(&models.Account{})
	suite.NewFactorySuite(suite.DB())

	NewRoutes(suite.App()).Setup(suite.DB())
}

// SetupTest will automatically refresh databases on every test. This can be overwritten in your test files.
func (suite *ControllerSuite) SetupTest() {
	suite.RefreshDB()
}

// TearDownSuite will close the underlying sqlite connection when the test suite is finished. This can be overwritten in your test files.
func (suite *ControllerSuite) TearDownSuite() {
	_ = suite.App().Shutdown()
	suite.ShutdownDB()
}

package ctrl

import (
	"github.com/netr/napi/examples/app/db/models"
	"github.com/netr/napi/trex"
	"github.com/stretchr/testify/suite"
	"testing"
)

type accountSuite struct {
	ControllerSuite
}

func (s *accountSuite) TestIndex_ExpectedBehavior() {
	a := s.CreateAccount()
	b := s.CreateAccount()

	trex.New(s).
		Get(s.Route("accounts.index"), nil).
		AssertOk().
		AssertDataCount(2).
		AssertJsonEqual("data[0].username", a.Username).
		AssertJsonEqual("data[1].username", b.Username)
}

func (s *accountSuite) TestStore_ExpectedBehavior() {
	pd := s.MakeUrlValues("username=testinghere&password=doingthisheresedsd")
	trex.New(s).
		Post(s.Route("accounts.store"), &pd, nil).
		AssertOk().
		AssertJsonEqual("data.username", "testinghere").
		AssertJsonEqual("data.id", float64(1))
}

func (s *accountSuite) TestStore_ShouldFail_PasswordTooShort() {
	pd := s.MakeUrlValues("username=testinghere&password=333")
	trex.New(s).
		Post(s.Route("accounts.store"), &pd, nil).
		AssertUnprocessable().
		AssertValidationErrors("password")
}

func (s *accountSuite) TestStore_ShouldFail_UsernameTooShort() {
	pd := s.MakeUrlValues("username=ad&password=asfsafasff")
	trex.New(s).
		Post(s.Route("accounts.store"), &pd, nil).
		AssertUnprocessable().
		AssertValidationErrors("username")
}

func (s *accountSuite) TestUpdate_ExpectedBehavior() {
	acc := s.CreateAccount()
	pd := s.MakeUrlValues("password=doingthisheresedsd")

	trex.New(s).
		Post(s.Route("accounts.update", acc.ID), &pd, nil).
		AssertOk().
		AssertJsonEqual("data.username", acc.Username).
		AssertJsonEqual("data.id", float64(1))

	s.AssertDatabaseHas(s.T(), &models.Account{Password: "doingthisheresedsd"})
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(accountSuite))
}

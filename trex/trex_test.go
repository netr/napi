package trex

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	mock := new(mockSuite)
	mock.app = fiber.New(fiber.Config{AppName: "testing app"})
	r := New(mock)
	assert.NotNil(t, r)
}

func Test_TestResponse_AssertOk(t *testing.T) {
	mock := newMockSuite(t)
	mock.app.Get("/test", func(ctx *fiber.Ctx) error {
		return nil
	})

	_ = New(mock).Get("/test", nil).AssertOk()
}

func Test_TestResponse_AssertStatus_ExpectedBehavior_WhenPathNotFound(t *testing.T) {
	mock := newMockSuite(t)
	mock.app.Get("/test3", func(ctx *fiber.Ctx) error {
		return nil
	})

	_ = New(mock).Get("/test", nil).AssertStatus(404)
}

func Test_TestResponse_Message_ExpectedBehavior(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	if tr.Message() != "success" {
		t.Fatal("message should have been `success`")
	}
}

func Test_TestResponse_AssertJsonEqual_ExpectedBehavior_WithMap(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: fiber.Map{
				"username": "testing",
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	tr.AssertJsonEqual("data.username", "testing")
}

func Test_TestResponse_AssertJsonEqual_ExpectedBehavior_WithMapArray(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: []fiber.Map{
				{"username": "testing"},
				{"username2": "testing2"},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	tr.AssertJsonEqual("data[1].username2", "testing2")
}

func Test_TestResponse_AssertJsonEqual_Bools(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: []fiber.Map{
				{"username": "testing"},
				{"username2": false},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	tr.AssertJsonEqual("data[1].username2", false)
}

func Test_TestResponse_AssertJsonEqual_NestedMaps(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: fiber.Map{
				"username": "testing",
				"nested":   fiber.Map{"developments": "here"},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	tr.AssertJsonEqual("data.nested.developments", "here")
}

func Test_TestResponse_AssertJsonEqual_WorkingWithArrayIdentification(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: []fiber.Map{
				{"username": "testing"},
				{"username2": false},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	tr.AssertJsonEqual("data[0].username", "testing")
}

func Test_TestResponse_AssertJsonEqual_ExpectedBehavior_WithIntegers(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: []fiber.Map{
				{"username": "testing"},
				{"username2": 2},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	tr.AssertJsonEqual("data[1].username2", float64(2))
}

func Test_TestResponse_AssertJsonEqual_ExpectedBehavior_WithBooleans(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: []fiber.Map{
				{"username": "testing"},
				{"username2": false},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	tr.AssertJsonEqual("data[1].username2", false)
}

func Test_TestResponse_DataToJson_ExpectedBehavior(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: []fiber.Map{
				{"username": "testing"},
				{"username2": false},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	assert.Equal(t, "testing", tr.DataToJson().GetIndex(0).Get("username").MustString())
}

func Test_TestResponse_ToJson_ExpectedBehavior(t *testing.T) {
	mock := newMockSuite(t)
	raw, err := json.Marshal(
		SuccessResponse{
			Message: "success",
			Data: []fiber.Map{
				{"username": "testing"},
				{"username2": false},
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	tr := mockTestResponseWithBytes(mock, raw)
	assert.Equal(t, "testing", tr.ToJson().Get("data").GetIndex(0).Get("username").MustString())
}

func mockTestResponseWithBytes(mock *mockSuite, raw []byte) *TestResponse {
	tr := &TestResponse{suite: mock, response: &http.Response{Body: io.NopCloser(bytes.NewReader(raw))}}
	return tr
}

func Test_TestResponse_AssertUnauthorized(t *testing.T) {
	mock := newMockSuite(t)
	mock.app.Get("/test", func(ctx *fiber.Ctx) error {
		_ = ctx.SendStatus(http.StatusUnauthorized)
		return nil
	})

	_ = New(mock).Get("/test", nil).AssertUnauthorized()
}

func Test_TestResponse_AssertUnprocessable(t *testing.T) {
	mock := newMockSuite(t)
	mock.app.Get("/test", func(ctx *fiber.Ctx) error {
		_ = ctx.SendStatus(http.StatusUnprocessableEntity)
		return nil
	})

	_ = New(mock).Get("/test", nil).AssertUnprocessable()
}

func Test_TestResponse_AssertValidationErrors(t *testing.T) {
	mock := newMockSuite(t)
	mock.app.Get("/test", func(ctx *fiber.Ctx) error {
		_ = ctx.SendStatus(http.StatusUnprocessableEntity)
		_ = ctx.JSON(fiber.Map{
			"message": "failed test",
			"errors": ErrorMap{
				"test1": "required",
				"test2": "required",
				"test3": "required",
				"test4": "should still work",
			},
		})
		return nil
	})

	_ = New(mock).Get("/test", nil).AssertValidationErrors("test1", "test2", "test3")
}

func Test_TestResponse_ParseFormErrors_ShouldWorkAsExpected(t *testing.T) {
	b := fiber.Map{
		"message": "failed test",
		"errors": ErrorMap{
			"test1": "required",
			"test2": "required",
			"test3": "required",
			"test4": "should still work",
		},
	}
	js, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	tr := &TestResponse{response: &http.Response{
		Body: io.NopCloser(bytes.NewReader(js)),
	}}
	fe, err := tr.ParseFormErrors()
	if err != nil {
		t.Fatal(err)
	}

	assert.Truef(t, fe.Contains("test1", "test2", "test3", "test4"), "wanted errors: [test1,test2,test3,test4], got: %v\n", fe.Fields())
}

func handleError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

type mockSuite struct {
	suite.Suite
	app *fiber.App
}

func newMockSuite(t *testing.T) *mockSuite {
	mock := new(mockSuite)
	mock.Suite.SetT(t)
	mock.app = fiber.New(fiber.Config{AppName: "testing app"})
	return mock
}

func (s *mockSuite) SetupTest()  {}
func (s *mockSuite) SetupSuite() {}
func (s *mockSuite) App() *fiber.App {
	return s.app
}

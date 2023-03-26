package resp

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSender_Success_ExpectedBehavior(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).Success("testing", nil)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSender_ExpectedBehavior_WhenReturningStructData(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		d := struct {
			Data string `json:"data"`
		}{Data: "testing"}

		err := New(c).Success("testing", d)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"data":{"data":"testing"}`)
}

func TestSender_Success_ExpectedBehavior_WhenReturningArrayData(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		dArr := []struct {
			Data string `json:"data"`
		}{
			{Data: "testing1"},
			{Data: "testing2"},
		}

		err := New(c).Success("testing", dArr)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"data":[{"data":"testing1"},{"data":"testing2"}]`)
}

func TestSender_Success_ShouldReturnEmptyObjectIfDataIsNil(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).Success("testing", nil)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"data":{}`)
}

func TestSender_Success_ShouldReturnEmptyArrayIfDataIsEmpty(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).Success("testing", make([]interface{}, 0))
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"data":[]`)
}

func TestSender_Error(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).Error("testing", errors.New("test error"))
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"error":"test error"`)
}

func TestSender_FormError(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).FormError("testing", ErrorBag{"field_name": "error message"})
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"errors":{"field_name":"error message"}`)
}

func TestSender_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).Unauthorized("testing", errors.New("test error"))
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"error":"test error"`)
}

func TestSender_BadRequest(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).BadRequest("testing", errors.New("test error"))
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"error":"test error"`)
}

func TestSender_NotFound(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		err := New(c).NotFound("testing", errors.New("test error"))
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	resp, err := newTestRequest(app, "GET", "/test", nil)
	handleError(t, err)
	body, err := responseToString(resp)
	handleError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Contains(t, body, "testing")
	assert.Contains(t, body, `"error":"test error"`)
}

func handleError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func responseToString(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func newTestRequest(app *fiber.App, method, target string, body io.Reader) (*http.Response, error) {
	req := httptest.NewRequest(method, target, body)
	return app.Test(req, 15000)
}

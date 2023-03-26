package napi

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func TestNewServer_DefaultFiberConfig(t *testing.T) {
	appName := "API"
	s := NewServer(DefaultFiberConfig(appName))

	if s.app.Config().AppName != appName {
		t.Fatalf("wanted %s, got: %s\n", appName, s.app.Config().AppName)
	}
}

func TestWithBaseMiddlewares_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithBaseMiddlewares(),
	)

	if s.app.HandlersCount() != 5 {
		t.Fatalf("wanted 5, got: %d\n", s.app.HandlersCount())
	}
}

func Test_NoMiddlewares_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
	)

	if s.app.HandlersCount() != 0 {
		t.Fatalf("wanted 0, got: %d\n", s.app.HandlersCount())
	}
}

func TestWithCatchAll_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithCatchAll(),
	)

	// same logic used in Run(). this tests the functionality. not the greatest, but it works.
	if !s.pathExists("*") {
		s.CatchAll()
	}

	body := testFailRequest(t, s)

	if !strings.Contains(string(body), "route '/fail' not found") {
		t.Fatalf("should have seen `route '/fail' not found` in:\n%s\n", string(body))
	}
}

func Test_NoCatchAll_ShouldNotGetNotFoundText(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
	)

	if s.pathExists("*") {
		t.Fatalf("should not have found catch all path: *")
	}

	body := testFailRequest(t, s)

	if strings.Contains(string(body), "Not Found") {
		t.Fatal("should not have found not found")
	}
}

func TestWithPort_ExpectedBehavior(t *testing.T) {
	s := NewServer(DefaultFiberConfig("test"), WithPort(1338))
	if s.port != 1338 {
		t.Fatalf("wanted port: 1338, got: %d\n", s.port)
	}
}

func TestWithPrometheus_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithPrometheus("test"),
	)

	// GET /metrics, HEAD /metrics, Prometheus Middleware
	if s.app.HandlersCount() != 3 {
		t.Fatalf("wanted 3, got: %d\n", s.app.HandlersCount())
	}

	if !s.methodAndPathExists("GET", "/metrics") {
		t.Fatalf("should have found route: GET /metrics")
	}
}

func TestWithLogger_ExpectedBehavior(t *testing.T) {
	var b bytes.Buffer
	s := NewServer(
		DefaultFiberConfig("test"),
		WithLogger(logger.Config{Output: &b}),
	)

	_ = testFailRequest(t, s)

	if b.Len() == 0 {
		t.Fatal("should have placed logs into the byte buffer")
	}
}

func TestWithLoggerOutput_ExpectedBehavior(t *testing.T) {
	var b bytes.Buffer
	s := NewServer(
		DefaultFiberConfig("test"),
		WithLoggerOutput(&b),
	)

	_ = testFailRequest(t, s)

	if b.Len() == 0 {
		t.Fatal("should have placed logs into the byte buffer")
	}
}

func TestWithDefaultLogger_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithDefaultLogger(),
	)

	if s.app.HandlersCount() != 1 {
		t.Fatalf("wanted 1 handler, got: %d\n", s.app.HandlersCount())
	}
}

func TestWithLoggerCallback_ExpectedBehavior(t *testing.T) {
	var b bytes.Buffer
	s := NewServer(
		DefaultFiberConfig("test"),
		WithLoggerDoneCallback(func(c *fiber.Ctx, logString []byte) {
			b.Write(logString)
		}),
	)

	_ = testFailRequest(t, s)

	if b.Len() == 0 {
		t.Fatal("should have placed logs into the byte buffer")
	}
}

func TestWithPprof_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithPprof(),
	)

	res := testPprofRequest(t, s, "/debug/pprof/")
	if !strings.Contains(string(res), ">/debug/pprof/<") {
		t.Fatal("should have found >/debug/pprof/< in pprof response")
	}
}

func TestWithPprof_WithPrefix_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithPprof("/test"),
	)

	res := testPprofRequest(t, s, "/test/debug/pprof/")
	if !strings.Contains(string(res), ">/debug/pprof/<") {
		t.Fatal("should have found >/debug/pprof/< in pprof response")
	}
}

func TestWithLimiter_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithLimiter(limiter.Config{}),
	)

	if s.app.HandlersCount() != 1 {
		t.Fatalf("wanted 1 handler, got: %d\n", s.app.HandlersCount())
	}
}

func TestWithDefaultLimiter_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithDefaultLimiter(),
	)

	if s.app.HandlersCount() != 1 {
		t.Fatalf("wanted 1 handler, got: %d\n", s.app.HandlersCount())
	}
}

func TestWithDefaultCORS_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithDefaultCORS(),
	)

	if s.app.HandlersCount() != 1 {
		t.Fatalf("wanted 1 handler, got: %d\n", s.app.HandlersCount())
	}
}

func TestWithCORS_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithCORS(cors.Config{}),
	)

	if s.app.HandlersCount() != 1 {
		t.Fatalf("wanted 1 handler, got: %d\n", s.app.HandlersCount())
	}
}

func TestWithCache_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithCache(defaultCacheConfig()),
	)

	if s.app.HandlersCount() != 1 {
		t.Fatalf("wanted 1 handler, got: %d\n", s.app.HandlersCount())
	}
}

func TestWithHealth_ExpectedBehavior(t *testing.T) {
	s := NewServer(
		DefaultFiberConfig("test"),
		WithHealth(),
	)

	if s.app.HandlersCount() != 2 {
		t.Fatalf("wanted 2 handlers (HEAD/GET), got: %d\n", s.app.HandlersCount())
	}
}

func ExampleNewServer() {
	_ = NewServer(
		DefaultFiberConfig("App Name"),
		WithCatchAll(), // Figure this out eventually.
	).
		Port(1338).
		UseBaseMiddlewares().
		UsePrometheus().
		UsePprof().
		UseHealth().
		UseDefaultCORS().
		UseCORS(cors.Config{}).
		UseDefaultLogger().
		UseLogger(logger.Config{}).
		UseLoggerOutput(&bytes.Buffer{}).
		UseLoggerDoneCallback(func(c *fiber.Ctx, logString []byte) {}).
		UseDefaultLimiter().
		UseLimiter(limiter.Config{}).
		UseDefaultCache()

	_ = NewServer(
		DefaultFiberConfig("App Name"),
		WithPort(1338),
		WithBaseMiddlewares(),
		WithHealth(),
		WithCatchAll(),
		WithDefaultCORS(),
		WithCORS(cors.Config{}),
		WithPrometheus(),
		WithPprof(),
		WithDefaultLogger(),
		WithLogger(logger.Config{}),
		WithLoggerOutput(&bytes.Buffer{}),
		WithLoggerDoneCallback(func(c *fiber.Ctx, logString []byte) {}),
		WithDefaultLimiter(),
		WithLimiter(limiter.Config{}),
		WithDefaultCache(),
	)
}

func testFailRequest(t *testing.T, s *Server) []byte {
	req := httptest.NewRequest("GET", "/fail", nil)
	resp, err := s.app.Test(req, 15000)
	if err != nil {
		t.Fatalf("testing http: %s\n", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %s\n", err)
	}

	return body
}

func testPprofRequest(t *testing.T, s *Server, path string) []byte {
	req := httptest.NewRequest("GET", path, nil)
	resp, err := s.app.Test(req, 15000)
	if err != nil {
		t.Fatalf("testing http: %s\n", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %s\n", err)
	}

	return body
}

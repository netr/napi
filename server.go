package napi

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/netr/napi/middleware"
)

// Server fiber app instance
type Server struct {
	app      *fiber.App
	catchAll bool
	port     int
}

// ServerOption type used for option pattern
type ServerOption func(*Server)

// NewServer Creates a new fiber instance with a fiber.Config and stackable options using the options pattern. DefaultFiberConfig() should be used in most cases.
func NewServer(fiberCfg fiber.Config, opts ...ServerOption) *Server {
	app := fiber.New(fiberCfg)
	s := &Server{
		app:  app,
		port: 1337,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithBaseMiddlewares a set of middlewares that are essentially plug and play. If these base middlewares need more flexibility, we can create individual options for them as we need them.
func WithBaseMiddlewares() ServerOption {
	return func(s *Server) {
		s.UseBaseMiddlewares()
	}
}

// WithPort set the web server port.
func WithPort(p int) ServerOption {
	return func(s *Server) {
		s.Port(p)
	}
}

// WithCatchAll sets up a simple catch all handler. This has to be a bool and used when Run() is called. If you set the catch all handler before the routes created by the application, everything will be caught. The bool removes this problem.
func WithCatchAll() ServerOption {
	return func(s *Server) {
		s.catchAll = true
	}
}

// WithPrometheus adds the middleware for prometheus.
func WithPrometheus(serviceName ...string) ServerOption {
	return func(s *Server) {
		s.UsePrometheus(serviceName...)
	}
}

// WithPprof adds the middleware for running pprof. This prefix will be added to the default path of "/debug/pprof/", for a resulting URL of: "/#endpoint#/debug/pprof/".
func WithPprof(endpoint ...string) ServerOption {
	return func(s *Server) {
		s.UsePprof(endpoint...)
	}
}

// WithLogger use the logger middleware with a custom logger.Config struct. Use this when you need full control of the logger. The other logger helpers are designed to be called on their own.
func WithLogger(cfg logger.Config) ServerOption {
	return func(s *Server) {
		s.UseLogger(cfg)
	}
}

// WithDefaultLogger use the logger middleware with the default configuration.
func WithDefaultLogger() ServerOption {
	return func(s *Server) {
		s.UseDefaultLogger()
	}
}

// WithLoggerOutput use the logger middleware with a custom output writer. Can use os.Stdout, os.File, bytes.Buffer, etc. This uses the default logger config. Meant to be used by itself.
func WithLoggerOutput(w io.Writer) ServerOption {
	return func(s *Server) {
		s.UseLoggerOutput(w)
	}
}

// WithLoggerDoneCallback use the logger middleware with a custom done callback after a log is written.
func WithLoggerDoneCallback(cb func(c *fiber.Ctx, logString []byte)) ServerOption {
	return func(s *Server) {
		s.UseLoggerDoneCallback(cb)
	}
}

// WithLimiter use the Limiter middleware with the a custom configuration.
func WithLimiter(cfg limiter.Config) ServerOption {
	return func(s *Server) {
		s.app.Use(limiter.New(cfg))
	}
}

// WithDefaultLimiter use the Limiter middleware with the default configuration.
func WithDefaultLimiter() ServerOption {
	return func(s *Server) {
		s.app.Use(limiter.New())
	}
}

// DefaultFiberConfig basic fiber configuration with write and read timeouts set under the hood to 30 seconds. Can expand this but might as well just create your own fiber.Config.
func DefaultFiberConfig(appName string) fiber.Config {
	return fiber.Config{
		AppName:      appName,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
}

// WithCORS use the CORS middleware with the a custom configuration.
func WithCORS(cfg cors.Config) ServerOption {
	return func(s *Server) {
		s.UseCORS(cfg)
	}
}

// WithDefaultCORS use the CORS middleware with the default configuration.
func WithDefaultCORS() ServerOption {
	return func(s *Server) {
		s.UseDefaultLogger()
	}
}

// WithHealth opens up a health ping endpoint to be used for uptime monitoring.
func WithHealth() ServerOption {
	return func(s *Server) {
		s.UseHealth()
	}
}

// WithCache uses cache control headers to set the cache control header to public and max age to 1 year.
func WithCache(cfg cache.Config) ServerOption {
	return func(s *Server) {
		s.UseCache(cfg)
	}
}

// WithDefaultCache uses cache control headers to set the cache control header to public and max age to 1 year.
func WithDefaultCache() ServerOption {
	return func(s *Server) {
		s.UseDefaultCache()
	}
}

// Run start the fiber server. Should always be called instead of s.App().Listen(). Uses a graceful shutdown mechanism from https://github.com/gofiber/recipes/blob/7a04f52833b70b97251d8a37893d2e0c599a8c15/graceful-shutdown/main.go
//
// Catch all needs to be called here to not interfere with routing being created by the application.
func (s *Server) Run() {
	if s.catchAll {
		if !s.pathExists("*") {
			s.CatchAll()
		}
	}

	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%d", s.port)); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Store channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	log.Println("Gracefully shutting down...")
	_ = s.app.Shutdown()

	log.Println("Running cleanup tasks...")
}

// CatchAll helper function to automatically catch bad urls
func (s *Server) CatchAll() fiber.Router {
	return s.app.All("*", func(c *fiber.Ctx) error {
		_ = c.SendStatus(http.StatusNotFound)
		return c.JSON(&fiber.Map{
			"message": fmt.Sprintf("route '%s' not found", c.OriginalURL()),
			"error":   "endpoint not found",
		})
	})
}

// App returns the underlying fiber app instance
func (s *Server) App() *fiber.App {
	return s.app
}

// UseLogger use the logger middleware with a custom logger.Config struct. Use this when you need full control of the logger. The other logger helpers are designed to be called on their own.
func (s *Server) UseLogger(cfg logger.Config) *Server {
	s.app.Use(logger.New(cfg))
	return s
}

// UseDefaultLogger use the logger middleware with the default configuration.
func (s *Server) UseDefaultLogger() *Server {
	s.app.Use(logger.New())
	return s
}

// UseLoggerOutput use the logger middleware with a custom output writer. Can use os.Stdout, os.File, bytes.Buffer, etc. This uses the default logger config. Meant to be used by itself.
func (s *Server) UseLoggerOutput(w io.Writer) *Server {
	cfg := defaultLoggerConfig()
	cfg.Output = w
	s.app.Use(logger.New(cfg))
	return s
}

// UseLoggerDoneCallback use the logger middleware with a custom done callback after a log is written.
func (s *Server) UseLoggerDoneCallback(cb func(c *fiber.Ctx, logString []byte)) *Server {
	cfg := defaultLoggerConfig()
	cfg.Done = cb
	s.app.Use(logger.New(cfg))
	return s
}

// UseBaseMiddlewares a set of middlewares that are essentially plug and play. If these base middlewares need more flexibility, we can create individual options for them as we need them.
func (s *Server) UseBaseMiddlewares() *Server {
	useBaseMiddlewares(s.app)
	return s
}

// Port helper function to set web server port.
func (s *Server) Port(p int) *Server {
	s.port = p
	return s
}

// UsePrometheus helper function to set prometheus middleware.
func (s *Server) UsePrometheus(serviceName ...string) *Server {
	sn := ToSnakeCase(s.app.Config().AppName)
	if len(serviceName) > 1 {
		sn = serviceName[0]
	}

	prometheus := middleware.NewPrometheus(sn)
	prometheus.RegisterAt(s.app, "/metrics")
	s.app.Use(prometheus.Middleware)
	return s
}

// UsePprof adds the middleware for running pprof. This prefix will be added to the default path of "/debug/pprof/", for a resulting URL of: "/#endpoint#/debug/pprof/".
func (s *Server) UsePprof(endpoint ...string) *Server {
	cfg := pprof.Config{Next: nil}
	if len(endpoint) > 0 {
		cfg.Prefix = endpoint[0]
	}

	s.app.Use(pprof.New(cfg))
	return s
}

// UseCORS use the CORS middleware with the a custom configuration.
func (s *Server) UseCORS(cfg cors.Config) *Server {
	s.app.Use(cors.New(cfg))
	return s
}

// UseDefaultCORS use the CORS middleware with the default configuration.
func (s *Server) UseDefaultCORS() *Server {
	s.app.Use(cors.New())
	return s
}

// UseLimiter use the Limiter middleware with the a custom configuration.
func (s *Server) UseLimiter(cfg limiter.Config) *Server {
	s.app.Use(limiter.New(cfg))
	return s
}

// UseDefaultLimiter use the Limiter middleware with the default configuration.
func (s *Server) UseDefaultLimiter() *Server {
	s.app.Use(limiter.New())
	return s
}

// UseHealth opens up a health ping endpoint to be used for uptime monitoring.
func (s *Server) UseHealth() *Server {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		_ = c.SendStatus(http.StatusOK)
		return c.JSON(&fiber.Map{
			"message": "OK",
		})
	})
	return s
}

// UseCache uses cache control headers to set the cache control header to public and max age to 1 year.
func (s *Server) UseCache(cfg cache.Config) *Server {
	s.app.Use(cache.New(cfg))
	return s
}

// UseDefaultCache uses cache control headers to set the cache control header to public and max age to 1 year.
func (s *Server) UseDefaultCache() *Server {
	s.app.Use(cache.New(defaultCacheConfig()))
	return s
}

// pathExists scans the app route stack for a matching path.
func (s *Server) pathExists(path string) bool {
	found := false
	for _, routes := range s.app.Stack() {
		for _, route := range routes {
			if route.Path == path {
				found = true
				break
			}
		}
	}

	return found
}

// methodAndPathExists scans the app route stack for a matching method and path.
func (s *Server) methodAndPathExists(method, path string) bool {
	found := false
	for _, routes := range s.app.Stack() {
		for _, route := range routes {
			if route.Path == path && route.Method == method {
				found = true
				break
			}
		}
	}

	return found
}

// useBaseMiddlewares plug and play middlewares for general use.
func useBaseMiddlewares(app *fiber.App) {
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(requestid.New())
	app.Use(recover.New())
}

// defaultLoggerConfig returns the default logger from the fiber docs.
func defaultLoggerConfig() logger.Config {
	return logger.Config{
		Next:         nil,
		Done:         nil,
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	}
}

// defaultCacheConfig is the default config
func defaultCacheConfig() cache.Config {
	return cache.Config{
		Next:         nil,
		Expiration:   1 * time.Minute,
		CacheControl: false,
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		},
		ExpirationGenerator:  nil,
		StoreResponseHeaders: false,
		Storage:              nil,
		MaxBytes:             0,
		Methods:              []string{fiber.MethodGet, fiber.MethodHead},
	}
}

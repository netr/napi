# NAPI
Netr API Framework. Opinionated API Framework for [Fiber](https://gofiber.io/). Took a lot of inspiration from the E2E testing in [Laravel](https://laravel.com).

This is primarily a framework for my own use, but I'm open to suggestions and PRs. I can write a better README and documentation if there's interest.

## Examples
- `/app` API Server with Gorm and Controller Testing

## Server

### Usage
```go
srv := NewServer(
    DefaultFiberConfig("App Name"),

    // Default is 1337 if you leave this option out
    WithPort(1338)
	
    // QoL
    WithBaseMiddlewares(),
    WithCatchAll(),
	
    // CORS
    WithDefaultCORS(),
    WithCORS(cors.Config{}),
	
    // Metrics
    WithPrometheus("app_name"),
	
    // Profiling
    WithPprof(),
	
    // Loggers
    WithDefaultLogger(),
    WithLogger(logger.Config{}),
    WithLoggerOutput(&bytes.Buffer{}),
    WithLoggerDoneCallback(func(c *fiber.Ctx, logString []byte) {}),

    // Limiter
    WithDefaultLimiter(),
    WithLimiter(limiter.Config{}),
)
srv.Run()
```

```go
srv := NewServer(
	    DefaultFiberConfig("App Name"),
	    WithCatchAll(), // Figure this out eventually.
	).
	Port(1338).
	UseBaseMiddlewares().
	UsePrometheus("app_name").
	UsePprof().
	UseDefaultCORS().
	UseCORS(cors.Config{}).
	UseDefaultLogger().
	UseLogger(logger.Config{}).
	UseLoggerOutput(&bytes.Buffer{}).
	UseLoggerDoneCallback(func(c *fiber.Ctx, logString []byte) {}).
	UseDefaultLimiter().
	UseLimiter(limiter.Config{})
    )
srv.Run()
```

## Testing controllers
```go 
type accountSuite struct {
    ControllerSuite
}

func (s *accountSuite) TestIndex_ExpectedBehavior() {
    a := s.CreateAccount()
    
    trex.New(s).
    Get(s.Route("accounts.index"), nil).
    AssertOk().
    AssertDataCount(2).
    AssertJsonEqual("data[0].username", a.Username).
    AssertJsonEqual("data[1].username", b.Username)
}
```

## rprint
Easily print your routes.

```go
// see Test_PrintPretty for more information
napi.NewRoutePrinter(fiberApp).PrintDebug()
napi.NewRoutePrinter(fiberApp).PrintPretty()

// Sugar
napi.PrintRoutes(fiberApp, isDebugMode)
```
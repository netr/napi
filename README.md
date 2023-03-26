# NAPI
The Netr API Framework is an opinionated API framework designed specifically for use with the Fiber web framework, a popular open-source web framework for Go programming language. This framework draws inspiration from Laravel's end-to-end testing capabilities, enabling developers to easily and effectively test their applications.

While this framework was initially developed for the personal use of its creator, it is open to suggestions and contributions from other developers via pull requests. As such, the developer is committed to continuously improving the framework and its accompanying documentation.

As of now, the Netr API Framework is still in the development stage. However, it is usable for small-scale projects and experimentation. Additionally, the development team is committed to ensuring that any new features added to the framework are designed in a way that does not break existing code and does not impact the underlying structure of the package.

-ChatGPT

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
    // Create account using factory [see /examples/db/models/suite.go & /factory]
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
package observability

import (
	fibertrace "github.com/DataDog/dd-trace-go/contrib/gofiber/fiber.v2/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/gofiber/fiber/v2"

	"github.com/cometagaming/ms-casino-go-v2/internal/config"
)

// InitTracer starts the Datadog APM tracer when DD_ENABLED=true and returns
// a stop function that must be deferred by the caller. When disabled the stop
// function is a no-op, so callers can always defer it unconditionally.
func InitTracer(cfg *config.DatadogConfig) func() {
	if !cfg.Enabled {
		return func() {}
	}
	tracer.Start(
		tracer.WithService(cfg.Service),
		tracer.WithEnv(cfg.Env),
		tracer.WithServiceVersion(cfg.Version),
	)
	return tracer.Stop
}

// FiberMiddleware returns the Datadog APM middleware for Fiber.
// When no tracer is running (DD_ENABLED=false) spans are no-ops with negligible overhead.
func FiberMiddleware() fiber.Handler {
	return fibertrace.Middleware()
}

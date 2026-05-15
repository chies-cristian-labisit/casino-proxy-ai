package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/handler"
	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/middleware"
)

// Setup registers all application routes on the Fiber app.
func Setup(app *fiber.App, health *handler.HealthHandler, customer *handler.CustomerHandler) {
	// Register global middleware in order:
	// 1. ErrorHandler first (catches panics and validation errors from all middlewares)
	// 2. TraceIdMiddleware (generates or validates traceId from headers)
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.TraceIdMiddleware)

	app.Get("/liveness", health.Liveness)
	app.Get("/readiness", health.Readiness)

	api := app.Group("/api/v2")
	api.Get("/customers/:idTx", customer.GetByIdTx)
}

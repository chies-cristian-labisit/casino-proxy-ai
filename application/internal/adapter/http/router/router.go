package router

import (
	"github.com/cometagaming/casino-proxy-ai/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

// Setup registers all application routes on the Fiber app.
func Setup(app *fiber.App, health *handler.HealthHandler, customer *handler.CustomerHandler) {
	app.Get("/liveness", health.Liveness)
	app.Get("/readiness", health.Readiness)

	api := app.Group("/api/v1")
	api.Get("/customers/:idTx", customer.GetByIdTx)
}

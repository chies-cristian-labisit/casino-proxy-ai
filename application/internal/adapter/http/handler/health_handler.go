package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// ReadinessChecker is a function that returns nil when the application is ready to serve traffic.
type ReadinessChecker func(ctx context.Context) error

// HealthHandler handles Kubernetes liveness and readiness probes.
type HealthHandler struct {
	readiness ReadinessChecker
}

// NewHealthHandler creates a HealthHandler with the given readiness checker.
func NewHealthHandler(readiness ReadinessChecker) *HealthHandler {
	return &HealthHandler{readiness: readiness}
}

// Liveness always returns 200 — the process is alive.
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "alive"})
}

// Readiness returns 200 when the readiness checker passes, 503 otherwise.
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	if err := h.readiness(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "not ready",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ready"})
}

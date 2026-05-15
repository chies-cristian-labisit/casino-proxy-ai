package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a standardized API error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"status"`
}

// ErrorHandler is a centralized error handler for all HTTP validation errors.
// It captures panics and validation errors and returns consistent error responses.
// Should be registered as middleware before other routes.
func ErrorHandler(c *fiber.Ctx) error {
	// Defer a recovery handler to catch panics.
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to string and respond with 500.
			errMsg := fmt.Sprintf("%v", r)
			c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "internal_server_error",
				Message: "An unexpected error occurred",
				Details: errMsg,
				Status:  fiber.StatusInternalServerError,
			})
		}
	}()

	// Continue to next handler.
	return c.Next()
}

package handler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/cometagaming/ms-casino-go-v2/internal/domain"
	"github.com/cometagaming/ms-casino-go-v2/internal/logging"
	"github.com/cometagaming/ms-casino-go-v2/internal/usecase"
)

// customerResponse represents the outgoing HTTP response.
type customerResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// CustomerHandler handles customer REST API requests.
type CustomerHandler struct {
	repo   usecase.CustomerRepository
	logger *logging.StructuredLogger
}

// NewCustomerHandler creates a CustomerHandler with the given repository.
func NewCustomerHandler(repo usecase.CustomerRepository) *CustomerHandler {
	logger := logging.New(slog.Default(), logging.LogFormatJSON)
	return &CustomerHandler{repo: repo, logger: logger}
}

// GetByIdTx handles GET /api/v1/customers/:idTx.
func (h *CustomerHandler) GetByIdTx(c *fiber.Ctx) error {
	ctx := c.UserContext()
	idTx := c.Params("idTx")

	// Log incoming request (INFO level) - use map for simple contextual data
	h.logger.Info(ctx, "customer lookup request received", map[string]interface{}{
		"code": idTx,
	})

	customer, err := h.repo.GetByCode(ctx, idTx)
	if err != nil {
		if errors.Is(err, domain.ErrCustomerNotFound) {
			// Not found is WARNING level (noteworthy but expected)
			h.logger.Warn(ctx, "customer not found", map[string]interface{}{"code": idTx})
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found"})
		}
		// Database errors are ERROR level (actual failures)
		h.logger.Error(ctx, "failed to fetch customer from repository", map[string]interface{}{
			"code":  idTx,
			"error": err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	response := customerResponse{
		ID:   customer.ID,
		Code: customer.Code,
		Name: customer.Name,
	}

	// Log successful response (INFO level)
	h.logger.Info(ctx, "customer lookup response", response)

	return c.Status(fiber.StatusOK).JSON(response)
}

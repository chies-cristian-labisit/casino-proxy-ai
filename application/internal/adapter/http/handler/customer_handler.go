package handler

import (
	"errors"

	"github.com/cometagaming/casino-proxy-ai/internal/domain"
	"github.com/cometagaming/casino-proxy-ai/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type customerResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// CustomerHandler handles customer REST API requests.
type CustomerHandler struct {
	repo usecase.CustomerRepository
}

// NewCustomerHandler creates a CustomerHandler with the given repository.
func NewCustomerHandler(repo usecase.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

// GetByIdTx handles GET /api/v1/customers/:idTx.
func (h *CustomerHandler) GetByIdTx(c *fiber.Ctx) error {
	idTx := c.Params("idTx")

	customer, err := h.repo.GetByCode(c.Context(), idTx)
	if err != nil {
		if errors.Is(err, domain.ErrCustomerNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(customerResponse{
		ID:   customer.ID,
		Code: customer.Code,
		Name: customer.Name,
	})
}

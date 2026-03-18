package handler

import (
	"net/http"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) Login(c fiber.Ctx) error {
	var input inputoutput.LoginInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "no data found",
		})
	}

	// check each field
	if input.Email == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "email is required"})
	}
	if input.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "password is required"})
	}

	output, err := h.loginUseCase.Execute(&input)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid credentials"})
	}

	if output == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "no user found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

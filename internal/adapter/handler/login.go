package handler

import (
	"net/http"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) Login(c fiber.Ctx) error {
	var input inputoutput.LoginInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "no data found",
		})
	}

	// check each field
	if input.Email == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "email is required",
		})
	}
	if input.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "password is required",
		})
	}

	output, err := h.loginUseCase.Execute(&input)
	if err != nil {
		return HandleError(c, err)
	}

	if output == nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

package handler

import (
	"fmt"
	"net/http"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterProfessor(c fiber.Ctx) error {
	var input inputoutput.RegisterProfessorInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "no data found",
		})
	}

	// check each field
	if input.Email == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "email is required",
		})
	}
	if input.Name == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "name is required",
		})
	}
	if input.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "password is required",
		})
	}
	if input.University == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "university is required",
		})
	}
	if input.Department == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "department is required",
		})
	}

	result, err := h.registerProfessorUsecase.Execute(&input)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "process fail",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": result.Message,
		"user_id": result.UserID,
	})
}

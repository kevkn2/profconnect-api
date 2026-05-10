package handler

import (
	"net/http"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// RegisterProfessor godoc
// @Summary Register professor user
// @Description Create a new professor account with email, name, password, university, and department
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body inputoutput.RegisterProfessorInput true "Professor registration credentials"
// @Success 200 {object} inputoutput.RegisterOutput
// @Failure 422 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/register/professor [post]
func (h *Handler) RegisterProfessor(c fiber.Ctx) error {
	var input inputoutput.RegisterProfessorInput
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
	if input.Name == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "name is required",
		})
	}
	if input.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "password is required",
		})
	}
	if input.University == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "university is required",
		})
	}
	if input.Department == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "department is required",
		})
	}

	result, err := h.registerProfessorUsecase.Execute(c.Context(), &input)
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": result.Message,
		"user_id": result.UserID,
	})
}

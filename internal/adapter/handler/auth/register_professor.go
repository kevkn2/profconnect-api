package auth_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
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
// @Failure 422 {object} shared_handler.ErrorResponse
// @Failure 400 {object} shared_handler.ErrorResponse
// @Router /api/auth/register/professor [post]
func (h *Handler) RegisterProfessor(c fiber.Ctx) error {
	var input inputoutput.RegisterProfessorInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "no data found",
		})
	}

	if input.Email == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "email is required",
		})
	}
	if input.Name == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "name is required",
		})
	}
	if input.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "password is required",
		})
	}
	if input.University == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "university is required",
		})
	}
	if input.Department == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "department is required",
		})
	}

	result, err := h.registerProfessorUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": result.Message,
		"user_id": result.UserID,
	})
}

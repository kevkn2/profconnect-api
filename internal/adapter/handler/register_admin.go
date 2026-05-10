package handler

import (
	"net/http"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// RegisterAdmin godoc
// @Summary Register admin user
// @Description Create a new admin account with email, name, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body inputoutput.RegisterInput true "Admin registration credentials"
// @Success 200 {object} inputoutput.RegisterOutput
// @Failure 422 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/register/admin [post]
func (h *Handler) RegisterAdmin(c fiber.Ctx) error {
	var input inputoutput.RegisterInput
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

	result, err := h.registerAdminUsecase.Execute(&input)
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": result.Message,
		"user_id": result.UserID,
	})
}

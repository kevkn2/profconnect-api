package auth_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body inputoutput.LoginInput true "Login credentials"
// @Success 200 {object} inputoutput.LoginOutput
// @Failure 422 {object} shared_handler.ErrorResponse
// @Failure 400 {object} shared_handler.ErrorResponse
// @Router /api/auth/login [post]
func (h *Handler) Login(c fiber.Ctx) error {
	var input inputoutput.LoginInput
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
	if input.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "password is required",
		})
	}

	output, err := h.loginUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}

	if output == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared_handler.ErrorResponse{
			Message: "invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "login successful",
		"accessToken":  output.AccessToken,
		"role":          output.Role,
		"refreshToken": output.RefreshToken,
		"type":          output.Type,
	})
}

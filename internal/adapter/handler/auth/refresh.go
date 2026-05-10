package auth_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// Refresh godoc
// @Summary Refresh access token
// @Description Exchange a valid refresh token for a new access + refresh token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body inputoutput.RefreshInput true "Refresh token"
// @Success 200 {object} inputoutput.RefreshOutput
// @Failure 400 {object} shared_handler.ErrorResponse
// @Failure 401 {object} shared_handler.ErrorResponse
// @Router /api/auth/refresh [post]
func (h *Handler) Refresh(c fiber.Ctx) error {
	var input inputoutput.RefreshInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "no data found",
		})
	}

	if input.RefreshToken == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "refresh_token is required",
		})
	}

	output, err := h.refreshUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":       "token refreshed",
		"accessToken":  output.AccessToken,
		"refreshToken": output.RefreshToken,
		"type":          output.Type,
	})
}

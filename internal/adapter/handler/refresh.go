package handler

import (
	"net/http"

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
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/refresh [post]
func (h *Handler) Refresh(c fiber.Ctx) error {
	var input inputoutput.RefreshInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "no data found",
		})
	}

	if input.RefreshToken == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "refresh_token is required",
		})
	}

	output, err := h.refreshUsecase.Execute(c.Context(), &input)
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":       "token refreshed",
		"token":         output.Token,
		"refresh_token": output.RefreshToken,
		"type":          output.Type,
	})
}

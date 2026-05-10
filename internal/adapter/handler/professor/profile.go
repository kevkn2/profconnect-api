package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// Profile godoc
// @Summary Get professor profile
// @Description Returns the authenticated professor's profile derived from JWT
// @Tags Professor
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.ProfessorProfileOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Router /api/professor/profile [get]
func (h *Handler) Profile(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	output, err := h.profileUsecase.Execute(c.Context(), &inputoutput.ProfileInput{UserID: userID})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(output)
}

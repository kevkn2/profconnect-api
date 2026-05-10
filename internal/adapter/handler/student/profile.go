package student_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// Profile godoc
// @Summary Get student profile
// @Description Returns the authenticated student's profile derived from JWT
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.StudentProfileOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Router /api/student/profile [get]
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

package student_handler

import (
	"net/http"
	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListMyApplications godoc
// @Summary List my applications
// @Description Student lists their own applications across all projects.
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.ListApplicationsOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Router /api/student/applications [get]
func (h *Handler) ListMyApplications(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	output, err := h.listMyApplicationsUsecase.Execute(c.Context(), &inputoutput.ListMyApplicationsInput{
		StudentUserID: userID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

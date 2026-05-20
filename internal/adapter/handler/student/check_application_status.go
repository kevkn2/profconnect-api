package student_handler

import (
	"net/http"
	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListApplicationsPerProject godoc
// @Summary List applications for a specific project
// @Description Check application status if pending or accepted.
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.CheckApplicationStatusOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Router /api/student/projects/{id}/applications [get]
func (h *Handler) ListApplicationsPerProject(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(http.StatusBadRequest).JSON(shared_handler.ErrorResponse{
			Message: "project id is required",
		})
	}

	output, err := h.listApplicationsPerIDUsecase.Execute(c.Context(), &inputoutput.CheckApplicationStatusInput{
		StudentUserID: userID,
		ProjectID:       projectID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

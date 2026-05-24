package student_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// LeaveProject godoc
// @Summary Leave a project
// @Description Student leaves a project they are an active member of. Frees the slot.
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} inputoutput.EmptyOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Router /api/student/projects/{id}/membership [delete]
func (h *Handler) LeaveProject(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "project id is required",
		})
	}

	output, err := h.leaveProjectUsecase.Execute(c.Context(), &inputoutput.LeaveProjectInput{
		StudentUserID: userID,
		ProjectID:     projectID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

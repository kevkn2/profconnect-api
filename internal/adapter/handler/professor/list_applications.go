package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListApplicationsByProject godoc
// @Summary List applications for a project
// @Description Professor lists applications for a project they own.
// @Tags Professor
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} inputoutput.ListProjectApplicationsByProjectOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Router /api/professor/projects/{id}/applications [get]
func (h *Handler) ListApplicationsByProject(c fiber.Ctx) error {
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

	output, err := h.listApplicationsByProjectUsecase.Execute(c.Context(), &inputoutput.ListApplicationsByProjectInput{
		ProfessorUserID: userID,
		ProjectID:       projectID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

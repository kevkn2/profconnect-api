package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListInvitationsByProject godoc
// @Summary List invitations sent for a project
// @Tags Professor
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} inputoutput.ListInvitationsOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Router /api/professor/projects/{id}/invitations [get]
func (h *Handler) ListInvitationsByProject(c fiber.Ctx) error {
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

	output, err := h.listInvitationsByProjectUsecase.Execute(c.Context(), &inputoutput.ListInvitationsByProjectInput{
		ProfessorUserID: userID,
		ProjectID:       projectID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

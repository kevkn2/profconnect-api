package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// SendInvitation godoc
// @Summary Invite a student to a project
// @Description Professor sends an invitation to a student for a project they own.
// @Tags Professor
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body inputoutput.SendInvitationInput true "Invitation payload"
// @Success 201 {object} inputoutput.ProjectInvitationOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Failure 409 {object} shared_handler.ErrorResponse
// @Router /api/professor/projects/{id}/invitations [post]
func (h *Handler) SendInvitation(c fiber.Ctx) error {
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

	var input inputoutput.SendInvitationInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "invalid request body",
		})
	}
	input.ProfessorUserID = userID
	input.ProjectID = projectID

	output, err := h.sendInvitationUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusCreated).JSON(output)
}

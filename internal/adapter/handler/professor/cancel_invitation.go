package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// CancelInvitation godoc
// @Summary Cancel a pending invitation
// @Tags Professor
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Param invId path string true "Invitation ID"
// @Success 200 {object} inputoutput.EmptyOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Failure 409 {object} shared_handler.ErrorResponse
// @Router /api/professor/projects/{id}/invitations/{invId} [delete]
func (h *Handler) CancelInvitation(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	projectID := c.Params("id")
	invID := c.Params("invId")
	if projectID == "" || invID == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "project id and invitation id are required",
		})
	}

	output, err := h.cancelInvitationUsecase.Execute(c.Context(), &inputoutput.CancelInvitationInput{
		ProfessorUserID: userID,
		ProjectID:       projectID,
		InvitationID:    invID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

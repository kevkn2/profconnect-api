package student_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// RespondInvitation godoc
// @Summary Accept or decline an invitation
// @Description Student accepts or declines an invitation. On accept, the student is added as a project member.
// @Tags Student
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param invId path string true "Invitation ID"
// @Param request body inputoutput.RespondInvitationInput true "Response payload"
// @Success 200 {object} inputoutput.ProjectInvitationOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Failure 409 {object} shared_handler.ErrorResponse
// @Router /api/student/invitations/{invId} [patch]
func (h *Handler) RespondInvitation(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	invID := c.Params("invId")
	if invID == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "invitation id is required",
		})
	}

	var input inputoutput.RespondInvitationInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "invalid request body",
		})
	}
	input.StudentUserID = userID
	input.InvitationID = invID

	output, err := h.respondInvitationUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

package student_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListMyInvitations godoc
// @Summary List my invitations
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.ListInvitationsOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Router /api/student/invitations [get]
func (h *Handler) ListMyInvitations(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	output, err := h.listMyInvitationsUsecase.Execute(c.Context(), &inputoutput.ListMyInvitationsInput{
		StudentUserID: userID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// RemoveMember godoc
// @Summary Remove a member from a project
// @Description Professor removes a member from a project they own. Frees the slot.
// @Tags Professor
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Param memberId path string true "Member ID"
// @Success 200 {object} inputoutput.EmptyOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Failure 409 {object} shared_handler.ErrorResponse
// @Router /api/professor/projects/{id}/members/{memberId} [delete]
func (h *Handler) RemoveMember(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	projectID := c.Params("id")
	memberID := c.Params("memberId")
	if projectID == "" || memberID == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "project id and member id are required",
		})
	}

	output, err := h.removeMemberUsecase.Execute(c.Context(), &inputoutput.RemoveMemberInput{
		ProfessorUserID: userID,
		ProjectID:       projectID,
		MemberID:        memberID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

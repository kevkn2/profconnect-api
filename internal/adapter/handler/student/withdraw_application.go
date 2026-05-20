package student_handler

import (
	"net/http"
	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// WithdrawApplication godoc
// @Summary Withdraw an application
// @Description Student withdraws a pending or rejected application.
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Param appId path string true "Application ID"
// @Success 200 {object} inputoutput.EmptyOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Failure 409 {object} shared_handler.ErrorResponse
// @Router /api/student/projects/{id}/applications/{appId} [delete]
func (h *Handler) WithdrawApplication(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	projectID := c.Params("id")
	appID := c.Params("appId")
	if projectID == "" || appID == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "project id and application id are required",
		})
	}

	output, err := h.withdrawApplicationUsecase.Execute(c.Context(), &inputoutput.WithdrawApplicationInput{
		StudentUserID: userID,
		ProjectID:     projectID,
		ApplicationID: appID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

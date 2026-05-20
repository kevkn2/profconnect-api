package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ReviewApplication godoc
// @Summary Approve or reject an application
// @Description Professor approves or rejects an application on a project they own.
// @Tags Professor
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param appId path string true "Application ID"
// @Param request body inputoutput.ReviewApplicationInput true "Review payload"
// @Success 200 {object} inputoutput.ProjectApplicationOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Failure 409 {object} shared_handler.ErrorResponse
// @Router /api/professor/projects/{id}/applications/{appId} [patch]
func (h *Handler) ReviewApplication(c fiber.Ctx) error {
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

	var input inputoutput.ReviewApplicationInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "no data found",
		})
	}
	input.ProfessorUserID = userID
	input.ProjectID = projectID
	input.ApplicationID = appID

	output, err := h.reviewApplicationUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

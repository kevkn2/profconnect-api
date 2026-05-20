package student_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ApplyProject godoc
// @Summary Apply to a project
// @Description Student applies to a project. Status is set to pending.
// @Tags Student
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body inputoutput.ApplyProjectInput false "Optional message"
// @Success 201 {object} inputoutput.ProjectApplicationOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Failure 409 {object} shared_handler.ErrorResponse
// @Router /api/student/projects/{id}/applications [post]
func (h *Handler) ApplyProject(c fiber.Ctx) error {
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

	var input inputoutput.ApplyProjectInput
	// Body is optional; ignore parsing errors and treat as empty message.
	_ = c.Bind().Body(&input)
	input.StudentUserID = userID
	input.ProjectID = projectID

	output, err := h.applyProjectUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusCreated).JSON(output)
}

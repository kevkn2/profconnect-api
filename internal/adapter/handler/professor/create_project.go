package professor_handler

import (
	"net/http"
	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// CreateProject godoc
// @Summary Create a project
// @Description Professor creates a new project. The professor_id is sourced from the JWT.
// @Tags Professor
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body inputoutput.CreateProjectInput true "Project payload"
// @Success 201 {object} inputoutput.ProjectOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 403 {object} shared_handler.ErrorResponse
// @Failure 422 {object} shared_handler.ErrorResponse
// @Router /api/professor/projects [post]
func (h *Handler) CreateProject(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	var input inputoutput.CreateProjectInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "no data found",
		})
	}
	input.ProfessorUserID = userID

	output, err := h.createProjectUsecase.Execute(c.Context(), &input)
	if err != nil {
		return shared_handler.HandleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(output)
}

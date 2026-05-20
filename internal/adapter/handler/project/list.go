package project_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListProjects godoc
// @Summary List all projects
// @Description Returns all projects available for students to browse.
// @Tags Project
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.ListProjectsOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Router /api/projects [get]
func (h *Handler) ListProjects(c fiber.Ctx) error {
	output, err := h.listProjectsUsecase.Execute(c.Context(), &inputoutput.ListProjectsInput{})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

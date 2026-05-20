package project_handler

import (
	"net/http"
	shared_handler "profconnect-api/internal/adapter/handler/shared"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// GetProject godoc
// @Summary Get a project by id
// @Tags Project
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} inputoutput.ProjectOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Router /api/projects/{id} [get]
func (h *Handler) GetProject(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "project id is required",
		})
	}

	output, err := h.getProjectUsecase.Execute(c.Context(), &inputoutput.GetProjectInput{ProjectID: id})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}
package project_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListMembers godoc
// @Summary List members of a project
// @Description Returns the roster of students who joined the project — via approved application or accepted invitation.
// @Tags Project
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} inputoutput.ListMembersOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Failure 404 {object} shared_handler.ErrorResponse
// @Router /api/projects/{id}/members [get]
func (h *Handler) ListMembers(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared_handler.ErrorResponse{
			Message: "project id is required",
		})
	}

	output, err := h.listMembersByProjectUsecase.Execute(c.Context(), &inputoutput.ListMembersByProjectInput{ProjectID: id})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

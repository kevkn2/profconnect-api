package student_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListMyProjects godoc
// @Summary List projects I am a member of
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.ListMyProjectsOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Router /api/student/projects [get]
func (h *Handler) ListMyProjects(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared_handler.ErrorResponse{
			Message: "missing user id in token",
		})
	}

	output, err := h.listMyProjectsUsecase.Execute(c.Context(), &inputoutput.ListMyProjectsInput{
		StudentUserID: userID,
	})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

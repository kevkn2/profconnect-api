package professor_handler

import (
	"net/http"

	shared_handler "profconnect-api/internal/adapter/handler/shared"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ListStudents godoc
// @Summary List all students
// @Description Returns the directory of students. Used by the frontend to pick whom to invite.
// @Tags Professor
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.ListStudentsOutput
// @Failure 401 {object} shared_handler.ErrorResponse
// @Router /api/professor/students [get]
func (h *Handler) ListStudents(c fiber.Ctx) error {
	output, err := h.listStudentsUsecase.Execute(c.Context(), &inputoutput.ListStudentsInput{})
	if err != nil {
		return shared_handler.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(output)
}

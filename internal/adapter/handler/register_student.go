package handler

import (
	"net/http"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// RegisterStudent godoc
// @Summary Register student user
// @Description Create a new student account with email, name, password, university, department, and research interests
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body inputoutput.RegisterStudentInput true "Student registration credentials"
// @Success 200 {object} inputoutput.RegisterOutput
// @Failure 422 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/register/student [post]
func (h *Handler) RegisterStudent(c fiber.Ctx) error {
	var input inputoutput.RegisterStudentInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "no data found",
		})
	}

	// check each field
	if input.Email == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "email is required",
		})
	}
	if input.Name == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "name is required",
		})
	}
	if input.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "password is required",
		})
	}
	if input.University == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "university is required",
		})
	}
	if input.Department == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "department is required",
		})
	}
	if input.ResearchInterests == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "research interests is required",
		})
	}

	result, err := h.registerStudentUsecase.Execute(&input)
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": result.Message,
		"user_id": result.UserID,
	})
}

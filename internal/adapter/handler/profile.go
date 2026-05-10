package handler

import (
	"net/http"

	"profconnect-api/internal/adapter/middleware"
	inputoutput "profconnect-api/internal/domain/input_output"

	"github.com/gofiber/fiber/v3"
)

// ProfessorProfile godoc
// @Summary Get professor profile
// @Description Returns the authenticated professor's profile derived from JWT
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.ProfessorProfileOutput
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/profile/professor [get]
func (h *Handler) ProfessorProfile(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
			Message: "missing user id in token",
		})
	}

	output, err := h.professorProfileUsecase.Execute(&inputoutput.ProfileInput{UserID: userID})
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(output)
}

// StudentProfile godoc
// @Summary Get student profile
// @Description Returns the authenticated student's profile derived from JWT
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} inputoutput.StudentProfileOutput
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/profile/student [get]
func (h *Handler) StudentProfile(c fiber.Ctx) error {
	userID, ok := c.Locals(middleware.LocalsUserID).(string)
	if !ok || userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
			Message: "missing user id in token",
		})
	}

	output, err := h.studentProfileUsecase.Execute(&inputoutput.ProfileInput{UserID: userID})
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(output)
}

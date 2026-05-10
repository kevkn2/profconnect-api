package student_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for student-only endpoints.
type Handler struct {
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput]
}

// New creates a new student handler.
func New(
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput],
) *Handler {
	return &Handler{
		profileUsecase: profileUsecase,
	}
}

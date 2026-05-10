package professor_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for professor-only endpoints.
type Handler struct {
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput]
}

// New creates a new professor handler.
func New(
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput],
) *Handler {
	return &Handler{
		profileUsecase: profileUsecase,
	}
}

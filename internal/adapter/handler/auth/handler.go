package auth_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for authentication endpoints.
type Handler struct {
	registerAdminUsecase     port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput]
	registerProfessorUsecase port.Usecase[inputoutput.RegisterProfessorInput, inputoutput.RegisterOutput]
	registerStudentUsecase   port.Usecase[inputoutput.RegisterStudentInput, inputoutput.RegisterOutput]
	loginUsecase             port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput]
	refreshUsecase           port.Usecase[inputoutput.RefreshInput, inputoutput.RefreshOutput]
}

// New creates a new auth handler.
func New(
	registerAdminUsecase port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput],
	registerProfessorUsecase port.Usecase[inputoutput.RegisterProfessorInput, inputoutput.RegisterOutput],
	registerStudentUsecase port.Usecase[inputoutput.RegisterStudentInput, inputoutput.RegisterOutput],
	loginUsecase port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput],
	refreshUsecase port.Usecase[inputoutput.RefreshInput, inputoutput.RefreshOutput],
) *Handler {
	return &Handler{
		registerAdminUsecase:     registerAdminUsecase,
		registerProfessorUsecase: registerProfessorUsecase,
		registerStudentUsecase:   registerStudentUsecase,
		loginUsecase:             loginUsecase,
		refreshUsecase:           refreshUsecase,
	}
}

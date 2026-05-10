package handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler contains all HTTP handlers
type Handler struct {
	registerAdminUsecase     port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput]
	registerProfessorUsecase port.Usecase[inputoutput.RegisterProfessorInput, inputoutput.RegisterOutput]
	registerStudentUsecase   port.Usecase[inputoutput.RegisterStudentInput, inputoutput.RegisterOutput]
	loginUseCase             port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput]
	professorProfileUsecase  port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput]
	studentProfileUsecase    port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput]
}

// NewHandler creates a new handler instance
func NewHandler(
	registerAdminUsecase port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput],
	registerProfessorUsecase port.Usecase[inputoutput.RegisterProfessorInput, inputoutput.RegisterOutput],
	registerStudentUsecase port.Usecase[inputoutput.RegisterStudentInput, inputoutput.RegisterOutput],
	loginUseCase port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput],
	professorProfileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput],
	studentProfileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput],
) *Handler {
	return &Handler{
		registerAdminUsecase:     registerAdminUsecase,
		registerProfessorUsecase: registerProfessorUsecase,
		registerStudentUsecase:   registerStudentUsecase,
		loginUseCase:             loginUseCase,
		professorProfileUsecase:  professorProfileUsecase,
		studentProfileUsecase:    studentProfileUsecase,
	}
}

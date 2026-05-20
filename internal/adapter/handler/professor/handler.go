package professor_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for professor-only endpoints.
type Handler struct {
	profileUsecase                   port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput]
	createProjectUsecase             port.Usecase[inputoutput.CreateProjectInput, inputoutput.ProjectOutput]
	listApplicationsByProjectUsecase port.Usecase[inputoutput.ListApplicationsByProjectInput, inputoutput.ListProjectApplicationsByProjectOutput]
	reviewApplicationUsecase         port.Usecase[inputoutput.ReviewApplicationInput, inputoutput.ProjectApplicationOutput]
}

// New creates a new professor handler.
func New(
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput],
	createProjectUsecase port.Usecase[inputoutput.CreateProjectInput, inputoutput.ProjectOutput],
	listApplicationsByProjectUsecase port.Usecase[inputoutput.ListApplicationsByProjectInput, inputoutput.ListProjectApplicationsByProjectOutput],
	reviewApplicationUsecase port.Usecase[inputoutput.ReviewApplicationInput, inputoutput.ProjectApplicationOutput],
) *Handler {
	return &Handler{
		profileUsecase:                   profileUsecase,
		createProjectUsecase:             createProjectUsecase,
		listApplicationsByProjectUsecase: listApplicationsByProjectUsecase,
		reviewApplicationUsecase:         reviewApplicationUsecase,
	}
}

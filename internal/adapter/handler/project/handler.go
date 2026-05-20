package project_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for shared project endpoints (browse/view).
type Handler struct {
	listProjectsUsecase port.Usecase[inputoutput.ListProjectsInput, inputoutput.ListProjectsOutput]
	getProjectUsecase   port.Usecase[inputoutput.GetProjectInput, inputoutput.ProjectOutput]
}

func New(
	listProjectsUsecase port.Usecase[inputoutput.ListProjectsInput, inputoutput.ListProjectsOutput],
	getProjectUsecase port.Usecase[inputoutput.GetProjectInput, inputoutput.ProjectOutput],
) *Handler {
	return &Handler{
		listProjectsUsecase: listProjectsUsecase,
		getProjectUsecase:   getProjectUsecase,
	}
}

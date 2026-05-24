package project_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for shared project endpoints (browse/view).
type Handler struct {
	listProjectsUsecase         port.Usecase[inputoutput.ListProjectsInput, inputoutput.ListProjectsOutput]
	getProjectUsecase           port.Usecase[inputoutput.GetProjectInput, inputoutput.ProjectOutput]
	listMembersByProjectUsecase port.Usecase[inputoutput.ListMembersByProjectInput, inputoutput.ListMembersOutput]
}

func New(
	listProjectsUsecase port.Usecase[inputoutput.ListProjectsInput, inputoutput.ListProjectsOutput],
	getProjectUsecase port.Usecase[inputoutput.GetProjectInput, inputoutput.ProjectOutput],
	listMembersByProjectUsecase port.Usecase[inputoutput.ListMembersByProjectInput, inputoutput.ListMembersOutput],
) *Handler {
	return &Handler{
		listProjectsUsecase:         listProjectsUsecase,
		getProjectUsecase:           getProjectUsecase,
		listMembersByProjectUsecase: listMembersByProjectUsecase,
	}
}

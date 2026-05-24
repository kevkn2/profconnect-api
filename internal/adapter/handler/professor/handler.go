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
	sendInvitationUsecase            port.Usecase[inputoutput.SendInvitationInput, inputoutput.ProjectInvitationOutput]
	listInvitationsByProjectUsecase  port.Usecase[inputoutput.ListInvitationsByProjectInput, inputoutput.ListInvitationsOutput]
	cancelInvitationUsecase          port.Usecase[inputoutput.CancelInvitationInput, inputoutput.EmptyOutput]
	removeMemberUsecase              port.Usecase[inputoutput.RemoveMemberInput, inputoutput.EmptyOutput]
	listStudentsUsecase              port.Usecase[inputoutput.ListStudentsInput, inputoutput.ListStudentsOutput]
}

// New creates a new professor handler.
func New(
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput],
	createProjectUsecase port.Usecase[inputoutput.CreateProjectInput, inputoutput.ProjectOutput],
	listApplicationsByProjectUsecase port.Usecase[inputoutput.ListApplicationsByProjectInput, inputoutput.ListProjectApplicationsByProjectOutput],
	reviewApplicationUsecase port.Usecase[inputoutput.ReviewApplicationInput, inputoutput.ProjectApplicationOutput],
	sendInvitationUsecase port.Usecase[inputoutput.SendInvitationInput, inputoutput.ProjectInvitationOutput],
	listInvitationsByProjectUsecase port.Usecase[inputoutput.ListInvitationsByProjectInput, inputoutput.ListInvitationsOutput],
	cancelInvitationUsecase port.Usecase[inputoutput.CancelInvitationInput, inputoutput.EmptyOutput],
	removeMemberUsecase port.Usecase[inputoutput.RemoveMemberInput, inputoutput.EmptyOutput],
	listStudentsUsecase port.Usecase[inputoutput.ListStudentsInput, inputoutput.ListStudentsOutput],
) *Handler {
	return &Handler{
		profileUsecase:                   profileUsecase,
		createProjectUsecase:             createProjectUsecase,
		listApplicationsByProjectUsecase: listApplicationsByProjectUsecase,
		reviewApplicationUsecase:         reviewApplicationUsecase,
		sendInvitationUsecase:            sendInvitationUsecase,
		listInvitationsByProjectUsecase:  listInvitationsByProjectUsecase,
		cancelInvitationUsecase:          cancelInvitationUsecase,
		removeMemberUsecase:              removeMemberUsecase,
		listStudentsUsecase:              listStudentsUsecase,
	}
}

package student_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for student-only endpoints.
type Handler struct {
	profileUsecase               port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput]
	applyProjectUsecase          port.Usecase[inputoutput.ApplyProjectInput, inputoutput.ProjectApplicationOutput]
	withdrawApplicationUsecase   port.Usecase[inputoutput.WithdrawApplicationInput, inputoutput.EmptyOutput]
	listMyApplicationsUsecase    port.Usecase[inputoutput.ListMyApplicationsInput, inputoutput.ListApplicationsOutput]
	listApplicationsPerIDUsecase port.Usecase[inputoutput.CheckApplicationStatusInput, inputoutput.CheckApplicationStatusOutput]
	listMyInvitationsUsecase     port.Usecase[inputoutput.ListMyInvitationsInput, inputoutput.ListInvitationsOutput]
	respondInvitationUsecase     port.Usecase[inputoutput.RespondInvitationInput, inputoutput.ProjectInvitationOutput]
	listMyProjectsUsecase        port.Usecase[inputoutput.ListMyProjectsInput, inputoutput.ListMyProjectsOutput]
	leaveProjectUsecase          port.Usecase[inputoutput.LeaveProjectInput, inputoutput.EmptyOutput]
}

// New creates a new student handler.
func New(
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput],
	applyProjectUsecase port.Usecase[inputoutput.ApplyProjectInput, inputoutput.ProjectApplicationOutput],
	withdrawApplicationUsecase port.Usecase[inputoutput.WithdrawApplicationInput, inputoutput.EmptyOutput],
	listMyApplicationsUsecase port.Usecase[inputoutput.ListMyApplicationsInput, inputoutput.ListApplicationsOutput],
	listApplicationsPerIDUsecase port.Usecase[inputoutput.CheckApplicationStatusInput, inputoutput.CheckApplicationStatusOutput],
	listMyInvitationsUsecase port.Usecase[inputoutput.ListMyInvitationsInput, inputoutput.ListInvitationsOutput],
	respondInvitationUsecase port.Usecase[inputoutput.RespondInvitationInput, inputoutput.ProjectInvitationOutput],
	listMyProjectsUsecase port.Usecase[inputoutput.ListMyProjectsInput, inputoutput.ListMyProjectsOutput],
	leaveProjectUsecase port.Usecase[inputoutput.LeaveProjectInput, inputoutput.EmptyOutput],
) *Handler {
	return &Handler{
		profileUsecase:               profileUsecase,
		applyProjectUsecase:          applyProjectUsecase,
		withdrawApplicationUsecase:   withdrawApplicationUsecase,
		listMyApplicationsUsecase:    listMyApplicationsUsecase,
		listApplicationsPerIDUsecase: listApplicationsPerIDUsecase,
		listMyInvitationsUsecase:     listMyInvitationsUsecase,
		respondInvitationUsecase:     respondInvitationUsecase,
		listMyProjectsUsecase:        listMyProjectsUsecase,
		leaveProjectUsecase:          leaveProjectUsecase,
	}
}

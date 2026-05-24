package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type cancelInvitationUsecase struct {
	projectRepository    port.ProjectRepository
	invitationRepository port.ProjectInvitationRepository
	professorRepository  port.ProfessorRepository
}

func NewCancelInvitationUsecase(
	projectRepository port.ProjectRepository,
	invitationRepository port.ProjectInvitationRepository,
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.CancelInvitationInput, inputoutput.EmptyOutput] {
	return &cancelInvitationUsecase{
		projectRepository:    projectRepository,
		invitationRepository: invitationRepository,
		professorRepository:  professorRepository,
	}
}

func (u *cancelInvitationUsecase) Execute(ctx context.Context, input *inputoutput.CancelInvitationInput) (*inputoutput.EmptyOutput, error) {
	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}

	professor, err := u.professorRepository.GetByUserID(ctx, input.ProfessorUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load professor", err)
	}
	if professor == nil {
		return nil, domain.NotFound("professor profile not found")
	}
	if project.Professor == nil || project.Professor.ID != professor.ID {
		return nil, domain.Forbidden("you are not the owner of this project")
	}

	invitation, err := u.invitationRepository.GetByID(ctx, input.InvitationID)
	if err != nil {
		return nil, domain.InternalErr("failed to load invitation", err)
	}
	if invitation == nil {
		return nil, domain.NotFound("invitation not found")
	}
	if invitation.Project.ID != project.ID {
		return nil, domain.NotFound("invitation does not belong to this project")
	}
	if invitation.Status != constants.InvitationStatusPending {
		return nil, domain.Conflict("only pending invitations can be cancelled")
	}

	if _, err := u.invitationRepository.UpdateStatus(ctx, invitation.ID, constants.InvitationStatusCancelled); err != nil {
		return nil, domain.InternalErr("failed to cancel invitation", err)
	}
	return &inputoutput.EmptyOutput{Message: "invitation cancelled"}, nil
}

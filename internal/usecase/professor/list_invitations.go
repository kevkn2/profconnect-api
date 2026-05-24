package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listInvitationsByProjectUsecase struct {
	projectRepository    port.ProjectRepository
	invitationRepository port.ProjectInvitationRepository
	professorRepository  port.ProfessorRepository
}

func NewListInvitationsByProjectUsecase(
	projectRepository port.ProjectRepository,
	invitationRepository port.ProjectInvitationRepository,
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.ListInvitationsByProjectInput, inputoutput.ListInvitationsOutput] {
	return &listInvitationsByProjectUsecase{
		projectRepository:    projectRepository,
		invitationRepository: invitationRepository,
		professorRepository:  professorRepository,
	}
}

func (u *listInvitationsByProjectUsecase) Execute(ctx context.Context, input *inputoutput.ListInvitationsByProjectInput) (*inputoutput.ListInvitationsOutput, error) {
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

	invitations, err := u.invitationRepository.ListByProject(ctx, project.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to list invitations", err)
	}

	out := make([]*inputoutput.ProjectInvitationOutput, 0, len(invitations))
	for _, inv := range invitations {
		out = append(out, invitationToOutput(inv, project, inv.Student))
	}
	return &inputoutput.ListInvitationsOutput{Invitations: out}, nil
}

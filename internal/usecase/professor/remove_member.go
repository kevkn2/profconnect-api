package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type removeMemberUsecase struct {
	projectRepository   port.ProjectRepository
	memberRepository    port.ProjectMemberRepository
	professorRepository port.ProfessorRepository
}

func NewRemoveMemberUsecase(
	projectRepository port.ProjectRepository,
	memberRepository port.ProjectMemberRepository,
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.RemoveMemberInput, inputoutput.EmptyOutput] {
	return &removeMemberUsecase{
		projectRepository:   projectRepository,
		memberRepository:    memberRepository,
		professorRepository: professorRepository,
	}
}

func (u *removeMemberUsecase) Execute(ctx context.Context, input *inputoutput.RemoveMemberInput) (*inputoutput.EmptyOutput, error) {
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

	member, err := u.memberRepository.GetByID(ctx, input.MemberID)
	if err != nil {
		return nil, domain.InternalErr("failed to load member", err)
	}
	if member == nil {
		return nil, domain.NotFound("member not found")
	}
	if member.Project.ID != project.ID {
		return nil, domain.NotFound("member does not belong to this project")
	}
	if member.Status != constants.MemberStatusActive {
		return nil, domain.Conflict("member is not active")
	}

	if _, err := u.memberRepository.UpdateStatus(ctx, member.ID, constants.MemberStatusRemoved); err != nil {
		return nil, domain.InternalErr("failed to remove member", err)
	}

	// Slot freed: if project was closed because it was full, reopen it.
	if project.Status == constants.ProjectStatusClosed {
		if err := u.projectRepository.UpdateStatus(ctx, project.ID, constants.ProjectStatusOpen); err != nil {
			return nil, domain.InternalErr("failed to reopen project", err)
		}
	}

	return &inputoutput.EmptyOutput{Message: "member removed"}, nil
}

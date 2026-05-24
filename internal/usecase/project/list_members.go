package project_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listMembersByProjectUsecase struct {
	projectRepository port.ProjectRepository
	memberRepository  port.ProjectMemberRepository
}

func NewListMembersByProjectUsecase(
	projectRepository port.ProjectRepository,
	memberRepository port.ProjectMemberRepository,
) port.Usecase[inputoutput.ListMembersByProjectInput, inputoutput.ListMembersOutput] {
	return &listMembersByProjectUsecase{
		projectRepository: projectRepository,
		memberRepository:  memberRepository,
	}
}

func (u *listMembersByProjectUsecase) Execute(ctx context.Context, input *inputoutput.ListMembersByProjectInput) (*inputoutput.ListMembersOutput, error) {
	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}

	members, err := u.memberRepository.ListByProject(ctx, project.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to list members", err)
	}

	out := make([]*inputoutput.ProjectMemberOutput, 0, len(members))
	for _, m := range members {
		out = append(out, MemberToOutput(m, project, m.Student))
	}
	return &inputoutput.ListMembersOutput{Members: out}, nil
}

// MemberToOutput converts a ProjectMember entity to its DTO. Exported so other
// usecases in sibling packages can reuse it.
func MemberToOutput(m *entities.ProjectMember, project *entities.Project, student *entities.Student) *inputoutput.ProjectMemberOutput {
	out := &inputoutput.ProjectMemberOutput{
		ID:     m.ID,
		Source: string(m.Source),
		Status: string(m.Status),
	}
	if !m.JoinedAt.IsZero() {
		out.JoinedAt = m.JoinedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	if m.LeftAt != nil {
		out.LeftAt = m.LeftAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	if m.SourceRefID != "" {
		out.SourceRefID = m.SourceRefID
	}
	if student != nil {
		brief := &inputoutput.ProjectStudentBrief{
			StudentID:         student.ID,
			University:        student.University,
			Department:        student.Department,
			ResearchInterests: student.ResearchInterests,
		}
		if student.User != nil {
			brief.UserID = student.User.ID
			brief.Name = student.User.Name
			brief.Email = student.User.Email
		}
		out.Student = brief
	}
	if project != nil {
		out.Project = &inputoutput.ProjectShortForApp{
			Title:       project.Title,
			Description: project.Description,
			Status:      string(project.Status),
		}
	}
	return out
}

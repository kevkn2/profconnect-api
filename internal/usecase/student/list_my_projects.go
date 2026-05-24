package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listMyProjectsUsecase struct {
	memberRepository  port.ProjectMemberRepository
	studentRepository port.StudentsRepository
}

func NewListMyProjectsUsecase(
	memberRepository port.ProjectMemberRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.ListMyProjectsInput, inputoutput.ListMyProjectsOutput] {
	return &listMyProjectsUsecase{
		memberRepository:  memberRepository,
		studentRepository: studentRepository,
	}
}

func (u *listMyProjectsUsecase) Execute(ctx context.Context, input *inputoutput.ListMyProjectsInput) (*inputoutput.ListMyProjectsOutput, error) {
	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	views, err := u.memberRepository.ListActiveByStudent(ctx, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to list memberships", err)
	}

	out := make([]*inputoutput.ProjectMemberOutput, 0, len(views))
	for _, v := range views {
		out = append(out, &inputoutput.ProjectMemberOutput{
			ID:     v.ID,
			Source: string(v.Source),
			Status: string(v.Status),
			Project: &inputoutput.ProjectShortForApp{
				Title:       v.ProjectTitle,
				Description: v.ProjectDescription,
				Status:      string(v.ProjectStatus),
			},
		})
	}
	return &inputoutput.ListMyProjectsOutput{Memberships: out}, nil
}

package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listMyInvitationsUsecase struct {
	invitationRepository port.ProjectInvitationRepository
	studentRepository    port.StudentsRepository
}

func NewListMyInvitationsUsecase(
	invitationRepository port.ProjectInvitationRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.ListMyInvitationsInput, inputoutput.ListInvitationsOutput] {
	return &listMyInvitationsUsecase{
		invitationRepository: invitationRepository,
		studentRepository:    studentRepository,
	}
}

func (u *listMyInvitationsUsecase) Execute(ctx context.Context, input *inputoutput.ListMyInvitationsInput) (*inputoutput.ListInvitationsOutput, error) {
	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	views, err := u.invitationRepository.ListByStudent(ctx, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to list invitations", err)
	}

	out := make([]*inputoutput.ProjectInvitationOutput, 0, len(views))
	for _, v := range views {
		out = append(out, &inputoutput.ProjectInvitationOutput{
			ID:      v.ID,
			Status:  string(v.Status),
			Message: v.Message,
			Project: &inputoutput.ProjectShortForApp{
				Title:       v.ProjectTitle,
				Description: v.ProjectDescription,
				Status:      string(v.ProjectStatus),
			},
		})
	}
	return &inputoutput.ListInvitationsOutput{Invitations: out}, nil
}

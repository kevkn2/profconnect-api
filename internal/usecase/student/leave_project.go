package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type leaveProjectUsecase struct {
	projectRepository port.ProjectRepository
	memberRepository  port.ProjectMemberRepository
	studentRepository port.StudentsRepository
}

func NewLeaveProjectUsecase(
	projectRepository port.ProjectRepository,
	memberRepository port.ProjectMemberRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.LeaveProjectInput, inputoutput.EmptyOutput] {
	return &leaveProjectUsecase{
		projectRepository: projectRepository,
		memberRepository:  memberRepository,
		studentRepository: studentRepository,
	}
}

func (u *leaveProjectUsecase) Execute(ctx context.Context, input *inputoutput.LeaveProjectInput) (*inputoutput.EmptyOutput, error) {
	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}

	member, err := u.memberRepository.GetActiveByProjectAndStudent(ctx, project.ID, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to load membership", err)
	}
	if member == nil {
		return nil, domain.NotFound("you are not an active member of this project")
	}

	if _, err := u.memberRepository.UpdateStatus(ctx, member.ID, constants.MemberStatusLeft); err != nil {
		return nil, domain.InternalErr("failed to leave project", err)
	}

	// Slot freed: if the project was closed because it was full, reopen it.
	if project.Status == constants.ProjectStatusClosed {
		if err := u.projectRepository.UpdateStatus(ctx, project.ID, constants.ProjectStatusOpen); err != nil {
			return nil, domain.InternalErr("failed to reopen project", err)
		}
	}

	return &inputoutput.EmptyOutput{Message: "left project"}, nil
}

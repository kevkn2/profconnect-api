package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type withdrawApplicationUsecase struct {
	applicationRepository port.ProjectApplicationRepository
	studentRepository     port.StudentsRepository
}

func NewWithdrawApplicationUsecase(
	applicationRepository port.ProjectApplicationRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.WithdrawApplicationInput, inputoutput.EmptyOutput] {
	return &withdrawApplicationUsecase{
		applicationRepository: applicationRepository,
		studentRepository:     studentRepository,
	}
}

func (u *withdrawApplicationUsecase) Execute(ctx context.Context, input *inputoutput.WithdrawApplicationInput) (*inputoutput.EmptyOutput, error) {
	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	application, err := u.applicationRepository.GetByID(ctx, input.ApplicationID)
	if err != nil {
		return nil, domain.InternalErr("failed to load application", err)
	}
	if application == nil {
		return nil, domain.NotFound("application not found")
	}
	if application.Project.ID != input.ProjectID {
		return nil, domain.NotFound("application does not belong to this project")
	}
	if application.Student == nil || application.Student.ID != student.ID {
		return nil, domain.Forbidden("you are not the owner of this application")
	}
	if application.Status == constants.ApplicationStatusApproved {
		return nil, domain.Conflict("approved applications cannot be withdrawn")
	}

	if err := u.applicationRepository.Delete(ctx, application.ID); err != nil {
		return nil, domain.InternalErr("failed to withdraw application", err)
	}

	return &inputoutput.EmptyOutput{Message: "application withdrawn"}, nil
}

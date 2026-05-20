package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listApplicationsPerIDUsecase struct {
	applicationRepository port.ProjectApplicationRepository
	studentRepository     port.StudentsRepository
}

func NewListApplicationsPerIDUsecase(
	applicationRepository port.ProjectApplicationRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.CheckApplicationStatusInput, inputoutput.CheckApplicationStatusOutput] {
	return &listApplicationsPerIDUsecase{
		applicationRepository: applicationRepository,
		studentRepository:     studentRepository,
	}
}

func (u *listApplicationsPerIDUsecase) Execute(ctx context.Context, input *inputoutput.CheckApplicationStatusInput) (*inputoutput.CheckApplicationStatusOutput, error) {
	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	exists, err := u.applicationRepository.CheckApplicationStatus(ctx, student.ID, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to check application", err)
	}

	return &inputoutput.CheckApplicationStatusOutput{Exists: exists}, nil
}

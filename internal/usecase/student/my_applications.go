package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listMyApplicationsUsecase struct {
	applicationRepository port.ProjectApplicationRepository
	studentRepository     port.StudentsRepository
}

func NewListMyApplicationsUsecase(
	applicationRepository port.ProjectApplicationRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.ListMyApplicationsInput, inputoutput.ListApplicationsOutput] {
	return &listMyApplicationsUsecase{
		applicationRepository: applicationRepository,
		studentRepository:     studentRepository,
	}
}

func (u *listMyApplicationsUsecase) Execute(ctx context.Context, input *inputoutput.ListMyApplicationsInput) (*inputoutput.ListApplicationsOutput, error) {
	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	views, err := u.applicationRepository.ListByStudent(ctx, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to list applications", err)
	}

	out := make([]*inputoutput.ProjectApplicationOutput, 0, len(views))
	for _, v := range views {
		out = append(out, &inputoutput.ProjectApplicationOutput{
			ID:        v.ID,
			Status:    string(v.Status),
			Message:   v.Message,
			Project: &inputoutput.ProjectShortForApp{
				Title:       v.ProjectTitle,
				Description: v.ProjectDescription,
				Status:      string(v.ProjectStatus),
			},
		})
	}

	return &inputoutput.ListApplicationsOutput{Applications: out}, nil
}

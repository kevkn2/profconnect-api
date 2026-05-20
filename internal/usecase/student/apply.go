package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type applyProjectUsecase struct {
	projectRepository     port.ProjectRepository
	applicationRepository port.ProjectApplicationRepository
	studentRepository     port.StudentsRepository
}

func NewApplyProjectUsecase(
	projectRepository port.ProjectRepository,
	applicationRepository port.ProjectApplicationRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.ApplyProjectInput, inputoutput.ProjectApplicationOutput] {
	return &applyProjectUsecase{
		projectRepository:     projectRepository,
		applicationRepository: applicationRepository,
		studentRepository:     studentRepository,
	}
}

func (u *applyProjectUsecase) Execute(ctx context.Context, input *inputoutput.ApplyProjectInput) (*inputoutput.ProjectApplicationOutput, error) {
	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}
	if project.Status != constants.ProjectStatusOpen {
		return nil, domain.Conflict("project is not accepting applications")
	}

	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	existing, err := u.applicationRepository.GetByProjectAndStudent(ctx, project.ID, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to check existing application", err)
	}
	if existing != nil {
		return nil, domain.Conflict("you have already applied to this project")
	}

	application := &entities.ProjectApplication{
		Project: project,
		Student:   student,
		Status:    constants.ApplicationStatusPending,
		Message:   input.Message,
	}

	created, err := u.applicationRepository.Create(ctx, application)
	if err != nil {
		return nil, domain.InternalErr("failed to create application", err)
	}

	return &inputoutput.ProjectApplicationOutput{
		ID:        created.ID,
		Project:   &inputoutput.ProjectShortForApp{
			Title: project.Title,
			Description: project.Description,
			Status: string(project.Status),
		},
		Student:   &inputoutput.ProjectStudentBrief{
			StudentID: created.Student.ID,
			UserID: created.Student.User.ID,
			Email: created.Student.User.Email,
			Name: created.Student.User.Name,
			University: created.Student.University,
			Department: created.Student.Department,
			ResearchInterests: created.Student.ResearchInterests,
		},
		Status:    string(created.Status),
		Message:   created.Message,
	}, nil
}
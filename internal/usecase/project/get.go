package project_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type getProjectUsecase struct {
	projectRepository port.ProjectRepository
}

func NewGetProjectUsecase(projectRepository port.ProjectRepository) port.Usecase[inputoutput.GetProjectInput, inputoutput.ProjectOutput] {
	return &getProjectUsecase{projectRepository: projectRepository}
}

func (u *getProjectUsecase) Execute(ctx context.Context, input *inputoutput.GetProjectInput) (*inputoutput.ProjectOutput, error) {
	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}
	return &inputoutput.ProjectOutput{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description,
		Slots:       project.Slots,
		Status:      string(project.Status),
		ProfessorID: project.Professor.ID,
		Professor: &inputoutput.ProjectProfessorBrief{
			UserID:     project.Professor.User.ID,
			Name:       project.Professor.User.Name,
			Email:      project.Professor.User.Email,
			University:  project.Professor.University,
			Department:  project.Professor.Department,
		},
	}, nil
}

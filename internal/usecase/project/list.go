package project_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listProjectsUsecase struct {
	projectRepository port.ProjectRepository
}

func NewListProjectsUsecase(projectRepository port.ProjectRepository) port.Usecase[inputoutput.ListProjectsInput, inputoutput.ListProjectsOutput] {
	return &listProjectsUsecase{projectRepository: projectRepository}
}

func (u *listProjectsUsecase) Execute(ctx context.Context, _ *inputoutput.ListProjectsInput) (*inputoutput.ListProjectsOutput, error) {
	projects, err := u.projectRepository.List(ctx)
	if err != nil {
		return nil, domain.InternalErr("failed to list projects", err)
	}

	out := make([]inputoutput.ProjectOutput, 0, len(projects))
	for _, p := range projects {
		out = append(out, inputoutput.ProjectOutput{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Slots:       p.Slots,
			Status:      string(p.Status),
			ProfessorID: p.Professor.ID,
			Professor: &inputoutput.ProjectProfessorBrief{
				UserID:     p.Professor.User.ID,
				Name:       p.Professor.User.Name,
				Email:      p.Professor.User.Email,
				University:  p.Professor.University,
				Department:  p.Professor.Department,
			},
		})
	}

	return &inputoutput.ListProjectsOutput{Projects: out}, nil
}

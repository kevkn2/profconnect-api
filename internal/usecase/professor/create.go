package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type createProjectUsecase struct {
	projectRepository   port.ProjectRepository
	professorRepository port.ProfessorRepository
}

func NewCreateProjectUsecase(
	projectRepository port.ProjectRepository,
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.CreateProjectInput, inputoutput.ProjectOutput] {
	return &createProjectUsecase{
		projectRepository:   projectRepository,
		professorRepository: professorRepository,
	}
}

func (u *createProjectUsecase) Execute(ctx context.Context, input *inputoutput.CreateProjectInput) (*inputoutput.ProjectOutput, error) {
	if input.Title == "" {
		return nil, domain.ValidationError("title is required")
	}
	if input.Description == "" {
		return nil, domain.ValidationError("description is required")
	}
	if input.Slots <= 0 {
		return nil, domain.ValidationError("slots must be greater than zero")
	}

	professor, err := u.professorRepository.GetByUserID(ctx, input.ProfessorUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load professor", err)
	}
	if professor == nil {
		return nil, domain.NotFound("professor profile not found")
	}

	project := &entities.Project{
		Professor:   professor,
		Title:       input.Title,
		Description: input.Description,
		Slots:       input.Slots,
		Status:      constants.ProjectStatusOpen,
	}

	created, err := u.projectRepository.Create(ctx, project)
	if err != nil {
		return nil, domain.InternalErr("failed to create project", err)
	}



	return &inputoutput.ProjectOutput{
		ID:          created.ID,
		Title:       created.Title,
		Description: created.Description,
		Slots:       created.Slots,
		Status:      string(created.Status),
	}, nil
}

func projectToOutput(p *entities.Project) *inputoutput.ProjectOutput {
	out := &inputoutput.ProjectOutput{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Slots:       p.Slots,
		Status:      string(p.Status),
	}
	if p.Professor != nil {
		out.ProfessorID = p.Professor.ID
		if p.Professor.User != nil {
			out.Professor = &inputoutput.ProjectProfessorBrief{
				UserID:     p.Professor.User.ID,
				Name:       p.Professor.User.Name,
				Email:      p.Professor.User.Email,
				University: p.Professor.University,
				Department: p.Professor.Department,
			}
		}
	}
	return out
}

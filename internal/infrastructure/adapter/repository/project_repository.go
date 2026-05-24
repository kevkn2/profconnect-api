package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"profconnect-api/internal/database/sqlc/generated"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	"profconnect-api/internal/domain/port"

	"github.com/google/uuid"
)

type projectRepository struct {
	queries *generated.Queries
}

func NewProjectRepository(queries *generated.Queries) port.ProjectRepository {
	return &projectRepository{queries: queries}
}

func (r *projectRepository) Create(ctx context.Context, project *entities.Project) (*entities.Project, error) {
	if project.Professor == nil {
		return nil, fmt.Errorf("professor is required")
	}
	professorID, err := uuid.Parse(project.Professor.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid professor id: %w", err)
	}

	row, err := r.queries.CreateProject(ctx, generated.CreateProjectParams{
		ProfessorID: professorID,
		Title:       project.Title,
		Description: project.Description,
		Slots:       int32(project.Slots),
		Status:      string(project.Status),
	})
	if err != nil {
		return nil, err
	}

	return &entities.Project{
		ID:          row.ID.String(),
		Professor:   project.Professor,
		Title:       row.Title,
		Description: row.Description,
		Slots:       int(row.Slots),
		Status:      constants.ProjectStatus(row.Status),
	}, nil
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*entities.Project, error) {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	row, err := r.queries.GetProjectByID(ctx, projectID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return rowToProject(row.ID, row.ProfessorID, row.Title, row.Description, row.Slots, row.Status,
		row.ProfessorUserID, row.ProfessorName, row.ProfessorEmail, row.ProfessorUniversity, row.ProfessorDepartment), nil
}

func (r *projectRepository) List(ctx context.Context) ([]*entities.Project, error) {
	rows, err := r.queries.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Project, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToProject(row.ID, row.ProfessorID, row.Title, row.Description, row.Slots, row.Status,
			row.ProfessorUserID, row.ProfessorName, row.ProfessorEmail, row.ProfessorUniversity, row.ProfessorDepartment))
	}
	return result, nil
}

func (r *projectRepository) ListByProfessor(ctx context.Context, professorID string) ([]*entities.Project, error) {
	profID, err := uuid.Parse(professorID)
	if err != nil {
		return nil, fmt.Errorf("invalid professor id: %w", err)
	}

	rows, err := r.queries.ListProjectsByProfessor(ctx, profID)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Project, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToProject(row.ID, row.ProfessorID, row.Title, row.Description, row.Slots, row.Status,
			row.ProfessorUserID, row.ProfessorName, row.ProfessorEmail, row.ProfessorUniversity, row.ProfessorDepartment))
	}
	return result, nil
}

func (r *projectRepository) UpdateStatus(ctx context.Context, id string, status constants.ProjectStatus) error {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid project id: %w", err)
	}

	return r.queries.UpdateProjectStatus(ctx, generated.UpdateProjectStatusParams{
		Status: string(status),
		ID:     projectID,
	})
}

func rowToProject(
	id, professorID uuid.UUID,
	title, description string,
	slots int32,
	status string,
	professorUserID uuid.UUID,
	professorName, professorEmail, professorUniversity, professorDepartment string,
) *entities.Project {
	return &entities.Project{
		ID: id.String(),
		Professor: &entities.Professor{
			ID: professorID.String(),
			User: &entities.User{
				ID:    professorUserID.String(),
				Name:  professorName,
				Email: professorEmail,
				Role:  constants.Professor,
			},
			University: professorUniversity,
			Department: professorDepartment,
		},
		Title:       title,
		Description: description,
		Slots:       int(slots),
		Status:      constants.ProjectStatus(status),
	}
}

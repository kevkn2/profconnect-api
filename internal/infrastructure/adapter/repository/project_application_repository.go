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

type projectApplicationRepository struct {
	queries *generated.Queries
}

func NewProjectApplicationRepository(queries *generated.Queries) port.ProjectApplicationRepository {
	return &projectApplicationRepository{queries: queries}
}

func (r *projectApplicationRepository) Create(ctx context.Context, app *entities.ProjectApplication) (*entities.ProjectApplication, error) {
	if app.Student == nil {
		return nil, fmt.Errorf("student is required")
	}
	projectID, err := uuid.Parse(app.Project.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}
	studentID, err := uuid.Parse(app.Student.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	row, err := r.queries.CreateProjectApplication(ctx, generated.CreateProjectApplicationParams{
		ProjectID: projectID,
		StudentID: studentID,
		Status:    string(app.Status),
		Message:   sql.NullString{String: app.Message, Valid: app.Message != ""},
	})
	if err != nil {
		return nil, err
	}

	return &entities.ProjectApplication{
		ID:        row.ID.String(),
		Project: &entities.Project{
			ID: row.ProjectID.String(),
		},
		Student:   app.Student,
		Status:    constants.ApplicationStatus(row.Status),
		Message:   row.Message.String,
	}, nil
}

func (r *projectApplicationRepository) GetByID(ctx context.Context, id string) (*entities.ProjectApplication, error) {
	appID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid application id: %w", err)
	}

	row, err := r.queries.GetProjectApplicationByID(ctx, appID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entities.ProjectApplication{
		ID:        row.ID.String(),
		Project: &entities.Project{
			ID:          row.ProjectID.String(),
		},
		Student: &entities.Student{
			ID: row.StudentID.String(),
			User: &entities.User{
				ID:    row.StudentUserID.String(),
				Name:  row.StudentName,
				Email: row.StudentEmail,
				Role:  constants.Student,
			},
			University:        row.StudentUniversity,
			Department:        row.StudentDepartment,
			ResearchInterests: row.StudentResearchInterests.String,
		},
		Status:  constants.ApplicationStatus(row.Status),
		Message: row.Message.String,
	}, nil
}

func (r *projectApplicationRepository) GetByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectApplication, error) {
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}
	studUUID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	row, err := r.queries.GetProjectApplicationByProjectAndStudent(ctx, generated.GetProjectApplicationByProjectAndStudentParams{
		ProjectID: projUUID,
		StudentID: studUUID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entities.ProjectApplication{
		ID:        row.ID.String(),
		Project: &entities.Project{
			ID:          row.ProjectID.String(),
		},
		Student:   &entities.Student{ID: row.StudentID.String()},
		Status:    constants.ApplicationStatus(row.Status),
		Message:   row.Message.String,
	}, nil
}

func (r *projectApplicationRepository) ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectApplication, error) {
	id, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	rows, err := r.queries.ListApplicationsByProject(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.ProjectApplication, 0, len(rows))
	for _, row := range rows {
		result = append(result, &entities.ProjectApplication{
			ID:        row.ID.String(),
			Project: &entities.Project{
				ID:          row.ProjectID.String(),
			},
			Student: &entities.Student{
				ID: row.StudentID.String(),
				User: &entities.User{
					ID:    row.StudentUserID.String(),
					Name:  row.StudentName,
					Email: row.StudentEmail,
					Role:  constants.Student,
				},
				University:        row.StudentUniversity,
				Department:        row.StudentDepartment,
				ResearchInterests: row.StudentResearchInterests.String,
			},
			Status:  constants.ApplicationStatus(row.Status),
			Message: row.Message.String,
		})
	}
	return result, nil
}

func (r *projectApplicationRepository) ListByStudent(ctx context.Context, studentID string) ([]*port.StudentApplicationView, error) {
	id, err := uuid.Parse(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	rows, err := r.queries.ListApplicationsByStudent(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]*port.StudentApplicationView, 0, len(rows))
	for _, row := range rows {
		result = append(result, &port.StudentApplicationView{
			ID:                 row.ID.String(),
			ProjectID:          row.ProjectID.String(),
			StudentID:          row.StudentID.String(),
			Status:             constants.ApplicationStatus(row.Status),
			Message:            row.Message.String,
			ProjectTitle:       row.ProjectTitle,
			ProjectDescription: row.ProjectDescription,
			ProjectStatus:      constants.ProjectStatus(row.ProjectStatus),
		})
	}
	return result, nil
}


func (r *projectApplicationRepository) CheckApplicationStatus(ctx context.Context, studentID, projectID string) (bool, error) {
	studentUUID, err := uuid.Parse(studentID)
	if err != nil {
		return false, fmt.Errorf("invalid student id: %w", err)
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return false, fmt.Errorf("invalid project id: %w", err)
	}

	rows, err := r.queries.CheckApplicationsByProjectAndStudent(ctx, generated.CheckApplicationsByProjectAndStudentParams{
		StudentID: studentUUID,
		ProjectID: projectUUID,
	})
	if err != nil {
		return false, err
	}

	return rows, nil
}

func (r *projectApplicationRepository) UpdateStatus(ctx context.Context, id string, status constants.ApplicationStatus) (*entities.ProjectApplication, error) {
	appID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid application id: %w", err)
	}

	row, err := r.queries.UpdateProjectApplicationStatus(ctx, generated.UpdateProjectApplicationStatusParams{
		Status: string(status),
		ID:     appID,
	})
	if err != nil {
		return nil, err
	}

	return &entities.ProjectApplication{
		ID:        row.ID.String(),
		Project: &entities.Project{
			ID:          row.ProjectID.String(),
		},
		Student:   &entities.Student{ID: row.StudentID.String()},
		Status:    constants.ApplicationStatus(row.Status),
		Message:   row.Message.String,
	}, nil
}

func (r *projectApplicationRepository) Delete(ctx context.Context, id string) error {
	appID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid application id: %w", err)
	}
	return r.queries.DeleteProjectApplication(ctx, appID)
}

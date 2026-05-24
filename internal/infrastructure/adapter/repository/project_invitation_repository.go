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

type projectInvitationRepository struct {
	queries *generated.Queries
}

func NewProjectInvitationRepository(queries *generated.Queries) port.ProjectInvitationRepository {
	return &projectInvitationRepository{queries: queries}
}

func (r *projectInvitationRepository) Create(ctx context.Context, inv *entities.ProjectInvitation) (*entities.ProjectInvitation, error) {
	if inv.Project == nil {
		return nil, fmt.Errorf("project is required")
	}
	if inv.Student == nil {
		return nil, fmt.Errorf("student is required")
	}

	projectID, err := uuid.Parse(inv.Project.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}
	studentID, err := uuid.Parse(inv.Student.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	status := inv.Status
	if status == "" {
		status = constants.InvitationStatusPending
	}

	row, err := r.queries.CreateProjectInvitation(ctx, generated.CreateProjectInvitationParams{
		ProjectID: projectID,
		StudentID: studentID,
		Status:    string(status),
		Message:   sql.NullString{String: inv.Message, Valid: inv.Message != ""},
	})
	if err != nil {
		return nil, err
	}

	return &entities.ProjectInvitation{
		ID:          row.ID.String(),
		Project:     &entities.Project{ID: row.ProjectID.String()},
		Student:     inv.Student,
		Status:      constants.InvitationStatus(row.Status),
		Message:     row.Message.String,
		RespondedAt: nullTimePtr(row.RespondedAt),
	}, nil
}

func (r *projectInvitationRepository) GetByID(ctx context.Context, id string) (*entities.ProjectInvitation, error) {
	invID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid invitation id: %w", err)
	}

	row, err := r.queries.GetProjectInvitationByID(ctx, invID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entities.ProjectInvitation{
		ID:      row.ID.String(),
		Project: &entities.Project{ID: row.ProjectID.String()},
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
		Status:      constants.InvitationStatus(row.Status),
		Message:     row.Message.String,
		RespondedAt: nullTimePtr(row.RespondedAt),
	}, nil
}

func (r *projectInvitationRepository) GetByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectInvitation, error) {
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}
	studUUID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	row, err := r.queries.GetInvitationByProjectAndStudent(ctx, generated.GetInvitationByProjectAndStudentParams{
		ProjectID: projUUID,
		StudentID: studUUID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entities.ProjectInvitation{
		ID:          row.ID.String(),
		Project:     &entities.Project{ID: row.ProjectID.String()},
		Student:     &entities.Student{ID: row.StudentID.String()},
		Status:      constants.InvitationStatus(row.Status),
		Message:     row.Message.String,
		RespondedAt: nullTimePtr(row.RespondedAt),
	}, nil
}

func (r *projectInvitationRepository) ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectInvitation, error) {
	id, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	rows, err := r.queries.ListInvitationsByProject(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.ProjectInvitation, 0, len(rows))
	for _, row := range rows {
		result = append(result, &entities.ProjectInvitation{
			ID:      row.ID.String(),
			Project: &entities.Project{ID: row.ProjectID.String()},
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
			Status:      constants.InvitationStatus(row.Status),
			Message:     row.Message.String,
			RespondedAt: nullTimePtr(row.RespondedAt),
		})
	}
	return result, nil
}

func (r *projectInvitationRepository) ListByStudent(ctx context.Context, studentID string) ([]*port.StudentInvitationView, error) {
	id, err := uuid.Parse(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	rows, err := r.queries.ListInvitationsByStudent(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]*port.StudentInvitationView, 0, len(rows))
	for _, row := range rows {
		result = append(result, &port.StudentInvitationView{
			ID:                 row.ID.String(),
			ProjectID:          row.ProjectID.String(),
			StudentID:          row.StudentID.String(),
			Status:             constants.InvitationStatus(row.Status),
			Message:            row.Message.String,
			ProjectTitle:       row.ProjectTitle,
			ProjectDescription: row.ProjectDescription,
			ProjectStatus:      constants.ProjectStatus(row.ProjectStatus),
		})
	}
	return result, nil
}

func (r *projectInvitationRepository) UpdateStatus(ctx context.Context, id string, status constants.InvitationStatus) (*entities.ProjectInvitation, error) {
	invID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid invitation id: %w", err)
	}

	row, err := r.queries.UpdateProjectInvitationStatus(ctx, generated.UpdateProjectInvitationStatusParams{
		Status: string(status),
		ID:     invID,
	})
	if err != nil {
		return nil, err
	}

	return &entities.ProjectInvitation{
		ID:          row.ID.String(),
		Project:     &entities.Project{ID: row.ProjectID.String()},
		Student:     &entities.Student{ID: row.StudentID.String()},
		Status:      constants.InvitationStatus(row.Status),
		Message:     row.Message.String,
		RespondedAt: nullTimePtr(row.RespondedAt),
	}, nil
}

func (r *projectInvitationRepository) Delete(ctx context.Context, id string) error {
	invID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid invitation id: %w", err)
	}
	return r.queries.DeleteProjectInvitation(ctx, invID)
}

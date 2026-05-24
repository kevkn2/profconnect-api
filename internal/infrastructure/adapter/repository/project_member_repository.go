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

type projectMemberRepository struct {
	queries *generated.Queries
}

func NewProjectMemberRepository(queries *generated.Queries) port.ProjectMemberRepository {
	return &projectMemberRepository{queries: queries}
}

func (r *projectMemberRepository) Create(ctx context.Context, member *entities.ProjectMember) (*entities.ProjectMember, error) {
	if member.Project == nil {
		return nil, fmt.Errorf("project is required")
	}
	if member.Student == nil {
		return nil, fmt.Errorf("student is required")
	}

	projectID, err := uuid.Parse(member.Project.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}
	studentID, err := uuid.Parse(member.Student.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	var sourceRef uuid.NullUUID
	if member.SourceRefID != "" {
		parsed, err := uuid.Parse(member.SourceRefID)
		if err != nil {
			return nil, fmt.Errorf("invalid source ref id: %w", err)
		}
		sourceRef = uuid.NullUUID{UUID: parsed, Valid: true}
	}

	status := member.Status
	if status == "" {
		status = constants.MemberStatusActive
	}

	row, err := r.queries.CreateProjectMember(ctx, generated.CreateProjectMemberParams{
		ProjectID:   projectID,
		StudentID:   studentID,
		Source:      string(member.Source),
		SourceRefID: sourceRef,
		Status:      string(status),
	})
	if err != nil {
		return nil, err
	}

	return &entities.ProjectMember{
		ID:          row.ID.String(),
		Project:     &entities.Project{ID: row.ProjectID.String()},
		Student:     member.Student,
		Source:      constants.MemberSource(row.Source),
		SourceRefID: nullUUIDString(row.SourceRefID),
		Status:      constants.MemberStatus(row.Status),
		JoinedAt:    row.JoinedAt.Time,
		LeftAt:      nullTimePtr(row.LeftAt),
	}, nil
}

func (r *projectMemberRepository) GetByID(ctx context.Context, id string) (*entities.ProjectMember, error) {
	memberID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid member id: %w", err)
	}

	row, err := r.queries.GetProjectMemberByID(ctx, memberID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entities.ProjectMember{
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
		Source:      constants.MemberSource(row.Source),
		SourceRefID: nullUUIDString(row.SourceRefID),
		Status:      constants.MemberStatus(row.Status),
		JoinedAt:    row.JoinedAt.Time,
		LeftAt:      nullTimePtr(row.LeftAt),
	}, nil
}

func (r *projectMemberRepository) GetActiveByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectMember, error) {
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}
	studUUID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	row, err := r.queries.GetActiveMemberByProjectAndStudent(ctx, generated.GetActiveMemberByProjectAndStudentParams{
		ProjectID: projUUID,
		StudentID: studUUID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entities.ProjectMember{
		ID:          row.ID.String(),
		Project:     &entities.Project{ID: row.ProjectID.String()},
		Student:     &entities.Student{ID: row.StudentID.String()},
		Source:      constants.MemberSource(row.Source),
		SourceRefID: nullUUIDString(row.SourceRefID),
		Status:      constants.MemberStatus(row.Status),
		JoinedAt:    row.JoinedAt.Time,
		LeftAt:      nullTimePtr(row.LeftAt),
	}, nil
}

func (r *projectMemberRepository) ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectMember, error) {
	id, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	rows, err := r.queries.ListMembersByProject(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.ProjectMember, 0, len(rows))
	for _, row := range rows {
		result = append(result, &entities.ProjectMember{
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
			Source:      constants.MemberSource(row.Source),
			SourceRefID: nullUUIDString(row.SourceRefID),
			Status:      constants.MemberStatus(row.Status),
			JoinedAt:    row.JoinedAt.Time,
			LeftAt:      nullTimePtr(row.LeftAt),
		})
	}
	return result, nil
}

func (r *projectMemberRepository) ListActiveByStudent(ctx context.Context, studentID string) ([]*port.StudentMembershipView, error) {
	id, err := uuid.Parse(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	rows, err := r.queries.ListActiveMembershipsByStudent(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]*port.StudentMembershipView, 0, len(rows))
	for _, row := range rows {
		result = append(result, &port.StudentMembershipView{
			ID:                 row.ID.String(),
			ProjectID:          row.ProjectID.String(),
			StudentID:          row.StudentID.String(),
			Source:             constants.MemberSource(row.Source),
			Status:             constants.MemberStatus(row.Status),
			ProjectTitle:       row.ProjectTitle,
			ProjectDescription: row.ProjectDescription,
			ProjectStatus:      constants.ProjectStatus(row.ProjectStatus),
			ProjectSlots:       int(row.ProjectSlots),
		})
	}
	return result, nil
}

func (r *projectMemberRepository) CountActive(ctx context.Context, projectID string) (int, error) {
	id, err := uuid.Parse(projectID)
	if err != nil {
		return 0, fmt.Errorf("invalid project id: %w", err)
	}

	count, err := r.queries.CountActiveMembersByProject(ctx, id)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *projectMemberRepository) UpdateStatus(ctx context.Context, id string, status constants.MemberStatus) (*entities.ProjectMember, error) {
	memberID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid member id: %w", err)
	}

	row, err := r.queries.UpdateProjectMemberStatus(ctx, generated.UpdateProjectMemberStatusParams{
		Status: string(status),
		ID:     memberID,
	})
	if err != nil {
		return nil, err
	}

	return &entities.ProjectMember{
		ID:          row.ID.String(),
		Project:     &entities.Project{ID: row.ProjectID.String()},
		Student:     &entities.Student{ID: row.StudentID.String()},
		Source:      constants.MemberSource(row.Source),
		SourceRefID: nullUUIDString(row.SourceRefID),
		Status:      constants.MemberStatus(row.Status),
		JoinedAt:    row.JoinedAt.Time,
		LeftAt:      nullTimePtr(row.LeftAt),
	}, nil
}

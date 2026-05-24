package repository

import (
	"context"
	"database/sql"
	"fmt"
	"profconnect-api/internal/database/sqlc/generated"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	"profconnect-api/internal/domain/port"

	"github.com/google/uuid"
)

type studentsRepository struct {
	queries *generated.Queries
}

func NewStudentsRepository(queries *generated.Queries) port.StudentsRepository {
	return &studentsRepository{
		queries: queries,
	}
}

// CreateStudent implements port.StudentsRepository.
func (s *studentsRepository) CreateStudent(ctx context.Context, student *entities.Student) (*entities.Student, error) {
	userID, err := uuid.Parse(student.User.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	resultStudent, err := s.queries.CreateStudent(ctx, generated.CreateStudentParams{
		UserID:            userID,
		University:        student.University,
		Department:        student.Department,
		ResearchInterests: sql.NullString{String: student.ResearchInterests, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}
	return &entities.Student{
		ID:                resultStudent.ID.String(),
		User:              student.User,
		University:        resultStudent.University,
		Department:        resultStudent.Department,
		ResearchInterests: resultStudent.ResearchInterests.String,
	}, nil
}

// GetStudentByID implements port.StudentsRepository.
func (s *studentsRepository) GetStudentByID(ctx context.Context, id string) (*entities.Student, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid student id: %w", err)
	}

	row, err := s.queries.GetStudentByID(ctx, parsedID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return &entities.Student{
		ID: row.ID.String(),
		User: &entities.User{
			ID:    row.UserID.String(),
			Name:  row.UserName,
			Email: row.UserEmail,
			Role:  constants.Student,
		},
		University:        row.University,
		Department:        row.Department,
		ResearchInterests: row.ResearchInterests.String,
	}, nil
}

// GetStudentByUserID implements port.StudentsRepository.
func (s *studentsRepository) GetStudentByUserID(ctx context.Context, userID string) (*entities.Student, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	resultStudent, err := s.queries.GetStudentByUserID(ctx, parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return &entities.Student{
		ID: resultStudent.ID.String(),
		User: &entities.User{
			ID:    resultStudent.UserID.String(),
			Name:  resultStudent.UserName,
			Email: resultStudent.UserEmail,
			Role: constants.Student,
		},
		University:        resultStudent.University,
		Department:        resultStudent.Department,
		ResearchInterests: resultStudent.ResearchInterests.String,
	}, nil
}

// ListStudents implements port.StudentsRepository.
func (s *studentsRepository) ListStudents(ctx context.Context) ([]*entities.Student, error) {
	rows, err := s.queries.ListStudents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	result := make([]*entities.Student, 0, len(rows))
	for _, row := range rows {
		result = append(result, &entities.Student{
			ID: row.ID.String(),
			User: &entities.User{
				ID:    row.UserID.String(),
				Name:  row.UserName,
				Email: row.UserEmail,
				Role:  constants.Student,
			},
			University:        row.University,
			Department:        row.Department,
			ResearchInterests: row.ResearchInterests.String,
		})
	}
	return result, nil
}

// UpdateStudent implements port.StudentsRepository.
func (s *studentsRepository) UpdateStudent(ctx context.Context, student *entities.Student) (*entities.Student, error) {
	resultStudent, err := s.queries.UpdateStudent(ctx, generated.UpdateStudentParams{
		University:        student.University,
		Department:        student.Department,
		ResearchInterests: sql.NullString{String: student.ResearchInterests, Valid: true},
		UserID:            uuid.MustParse(student.User.ID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}

	return &entities.Student{
		ID: resultStudent.ID.String(),
		User: &entities.User{
			ID:    resultStudent.UserID.String(),
			Name:  resultStudent.Name,
			Email: resultStudent.Email,
			Role:  constants.Student,
		},
		University:        resultStudent.University,
		Department:        resultStudent.Department,
		ResearchInterests: resultStudent.ResearchInterests.String,
	}, nil
}

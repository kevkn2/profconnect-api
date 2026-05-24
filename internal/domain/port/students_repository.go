package port

import (
	"context"
	"profconnect-api/internal/domain/entities"
)

type StudentsRepository interface {
	CreateStudent(ctx context.Context, student *entities.Student) (*entities.Student, error)
	GetStudentByID(ctx context.Context, id string) (*entities.Student, error)
	GetStudentByUserID(ctx context.Context, userID string) (*entities.Student, error)
	ListStudents(ctx context.Context) ([]*entities.Student, error)
	UpdateStudent(ctx context.Context, student *entities.Student) (*entities.Student, error)
}

package port

import (
	"context"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *entities.Project) (*entities.Project, error)
	GetByID(ctx context.Context, id string) (*entities.Project, error)
	List(ctx context.Context) ([]*entities.Project, error)
	ListByProfessor(ctx context.Context, professorID string) ([]*entities.Project, error)
	UpdateStatus(ctx context.Context, id string, status constants.ProjectStatus) error
	CountApprovedApplications(ctx context.Context, projectID string) (int, error)
}

type ProjectApplicationRepository interface {
	Create(ctx context.Context, application *entities.ProjectApplication) (*entities.ProjectApplication, error)
	GetByID(ctx context.Context, id string) (*entities.ProjectApplication, error)
	GetByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectApplication, error)
	ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectApplication, error)
	ListByStudent(ctx context.Context, studentID string) ([]*StudentApplicationView, error)
	CheckApplicationStatus(ctx context.Context, studentID, projectID string) (bool, error)
	UpdateStatus(ctx context.Context, id string, status constants.ApplicationStatus) (*entities.ProjectApplication, error)
	Delete(ctx context.Context, id string) error
}

// StudentApplicationView is a denormalized view a student sees of their own application,
// including project title/description without requiring a separate fetch.
type StudentApplicationView struct {
	ID                 string
	ProjectID          string
	StudentID          string
	Status             constants.ApplicationStatus
	Message            string
	ProjectTitle       string
	ProjectDescription string
	ProjectStatus      constants.ProjectStatus
}

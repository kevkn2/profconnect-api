package port

import (
	"context"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
)

type ProjectInvitationRepository interface {
	Create(ctx context.Context, invitation *entities.ProjectInvitation) (*entities.ProjectInvitation, error)
	GetByID(ctx context.Context, id string) (*entities.ProjectInvitation, error)
	GetByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectInvitation, error)
	ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectInvitation, error)
	ListByStudent(ctx context.Context, studentID string) ([]*StudentInvitationView, error)
	UpdateStatus(ctx context.Context, id string, status constants.InvitationStatus) (*entities.ProjectInvitation, error)
	Delete(ctx context.Context, id string) error
}

// StudentInvitationView is a denormalized view of an invitation as a student sees it,
// including project title/description without requiring a separate fetch.
type StudentInvitationView struct {
	ID                 string
	ProjectID          string
	StudentID          string
	Status             constants.InvitationStatus
	Message            string
	ProjectTitle       string
	ProjectDescription string
	ProjectStatus      constants.ProjectStatus
}

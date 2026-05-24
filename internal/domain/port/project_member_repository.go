package port

import (
	"context"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
)

type ProjectMemberRepository interface {
	Create(ctx context.Context, member *entities.ProjectMember) (*entities.ProjectMember, error)
	GetByID(ctx context.Context, id string) (*entities.ProjectMember, error)
	GetActiveByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectMember, error)
	ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectMember, error)
	ListActiveByStudent(ctx context.Context, studentID string) ([]*StudentMembershipView, error)
	CountActive(ctx context.Context, projectID string) (int, error)
	UpdateStatus(ctx context.Context, id string, status constants.MemberStatus) (*entities.ProjectMember, error)
}

// StudentMembershipView is a denormalized view a student sees of their own membership,
// including project title/description without requiring a separate fetch.
type StudentMembershipView struct {
	ID                 string
	ProjectID          string
	StudentID          string
	Source             constants.MemberSource
	Status             constants.MemberStatus
	ProjectTitle       string
	ProjectDescription string
	ProjectStatus      constants.ProjectStatus
	ProjectSlots       int
}

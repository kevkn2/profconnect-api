package port

import (
	"context"
	"profconnect-api/internal/domain/entities"
)

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
	Delete(ctx context.Context, id string) error
}

package port

import "profconnect-api/internal/domain/entities"

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	GetByEmail(email string) (*entities.User, error)
	Create(user *entities.User) (*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	Delete(id string) error
}

package repository

import (
	"context"
	"fmt"

	"profconnect-api/internal/database/sqlc/generated"
	"profconnect-api/internal/domain/entities"
	"profconnect-api/internal/domain/port"

	"github.com/google/uuid"
)

// UserRepository implements the UserRepository port interface
type userRepository struct {
	queries *generated.Queries
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(queries *generated.Queries) port.UserRepository {
	return &userRepository{
		queries: queries,
	}
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	user, err := r.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, nil // Return nil for not found to match port contract
	}

	return mapGetUserByEmailRowToEntity(user), nil
}

// Create inserts a new user
func (r *userRepository) Create(user *entities.User) (*entities.User, error) {
	created, err := r.queries.CreateUser(context.Background(), generated.CreateUserParams{
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user.ID = created.ID.String()
	return user, nil
}

// Update modifies an existing user
func (r *userRepository) Update(user *entities.User) (*entities.User, error) {
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	_, err = r.queries.UpdateUser(context.Background(), generated.UpdateUserParams{
		ID:             userID,
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

// Delete removes a user
func (r *userRepository) Delete(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	err = r.queries.DeleteUserByID(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

// mapGetUserByEmailRowToEntity converts GetUserByEmailRow to domain entity
func mapGetUserByEmailRowToEntity(u generated.GetUserByEmailRow) *entities.User {
	return &entities.User{
		ID:             u.ID.String(),
		Name:           u.Name,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
	}
}

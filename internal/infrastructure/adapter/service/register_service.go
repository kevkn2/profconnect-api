package service

import (
	"context"
	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	profconnect_utils "profconnect-api/internal/infrastructure/pkg/utils"
)

type registerService struct {
	userRepository port.UserRepository
}

func NewRegisterService(userRepository port.UserRepository) port.Service[inputoutput.RegisterInput, entities.User] {
	return &registerService{
		userRepository: userRepository,
	}
}

// Execute implements port.Service.
func (r *registerService) Execute(ctx context.Context, input *inputoutput.RegisterInput) (*entities.User, error) {
	// Check if user with this email already exists
	userExist, err := r.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.InternalErr("failed to check existing user", err)
	}

	if userExist != nil {
		return nil, domain.Conflict("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := profconnect_utils.HashPassword(input.Password)
	if err != nil {
		return nil, domain.InternalErr("failed to hash password", err)
	}

	// Create new user
	newUser := &entities.User{
		Name:           input.Name,
		Email:          input.Email,
		HashedPassword: hashedPassword,
		Role:           constants.Roles(input.Role),
	}

	// Save user to database
	createdUser, err := r.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, domain.InternalErr("failed to create user", err)
	}

	return createdUser, nil
}

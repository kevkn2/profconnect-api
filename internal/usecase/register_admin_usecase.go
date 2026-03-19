package usecase

import (
	"context"
	"errors"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	profconnect_utils "profconnect-api/internal/infrastructure/pkg/utils"
)

type registerAdminUsecase struct {
	userRepository port.UserRepository
}

// NewRegisterAdminUsecase creates a new instance of RegisterAdminUsecase
func NewRegisterAdminUsecase(userRepository port.UserRepository) port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput] {
	return &registerAdminUsecase{
		userRepository: userRepository,
	}
}

// Execute implements port.Usecase
func (r *registerAdminUsecase) Execute(input *inputoutput.RegisterInput) (*inputoutput.RegisterOutput, error) {
	ctx := context.Background()
	email := input.Email
	name := input.Name
	password := input.Password

	// Check if user with this email already exists
	userExist, err := r.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if userExist != nil {
		return nil, errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := profconnect_utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user
	newUser := &entities.User{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
		Role:           string(constants.Admin),
	}

	// Save user to database
	createdUser, err := r.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &inputoutput.RegisterOutput{
		Message: "User registered successfully",
		UserID:  createdUser.ID,
	}, nil
}

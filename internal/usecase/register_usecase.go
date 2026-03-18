package usecase

import (
	"errors"

	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	profconnect_utils "profconnect-api/internal/infrastructure/pkg/utils"
)

type registerUsecase struct {
	userRepository port.UserRepository
}

// NewRegisterUsecase creates a new instance of RegisterUsecase
func NewRegisterUsecase(userRepository port.UserRepository) port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput] {
	return &registerUsecase{
		userRepository: userRepository,
	}
}

// Execute implements port.Usecase
func (r *registerUsecase) Execute(input *inputoutput.RegisterInput) (*inputoutput.RegisterOutput, error) {
	email := input.Email
	name := input.Name
	password := input.Password

	// Check if user with this email already exists
	userExist, err := r.userRepository.GetByEmail(email)
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
	}

	// Save user to database
	createdUser, err := r.userRepository.Create(newUser)
	if err != nil {
		return nil, err
	}

	return &inputoutput.RegisterOutput{
		Message: "User registered successfully",
		UserID:  createdUser.ID,
	}, nil
}

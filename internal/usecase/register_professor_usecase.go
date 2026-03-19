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

type registerProfessorUsecase struct {
	userRepository      port.UserRepository
	professorRepository port.ProfessorRepository
}

// NewRegisterProfessorUsecase creates a new instance of RegisterProfessorUsecase
func NewRegisterProfessorUsecase(
	userRepository port.UserRepository,
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.RegisterProfessorInput, inputoutput.RegisterOutput] {
	return &registerProfessorUsecase{
		userRepository:      userRepository,
		professorRepository: professorRepository,
	}
}

// Execute implements port.Usecase
func (r *registerProfessorUsecase) Execute(input *inputoutput.RegisterProfessorInput) (*inputoutput.RegisterOutput, error) {
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

	// save the professor data to database
	newProfessor := &entities.Professor{
		User:       createdUser,
		University: input.University,
		Department: input.Department,
	}

	// Save professor to database
	_, err = r.professorRepository.Create(ctx, newProfessor)
	if err != nil {
		return nil, err
	}

	return &inputoutput.RegisterOutput{
		Message: "User registered successfully",
		UserID:  createdUser.ID,
	}, nil
}

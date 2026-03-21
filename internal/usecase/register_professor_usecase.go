package usecase

import (
	"context"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type registerProfessorUsecase struct {
	registerService     port.Service[inputoutput.RegisterInput, entities.User]
	professorRepository port.ProfessorRepository
}

// NewRegisterProfessorUsecase creates a new instance of RegisterProfessorUsecase
func NewRegisterProfessorUsecase(
	registerService port.Service[inputoutput.RegisterInput, entities.User],
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.RegisterProfessorInput, inputoutput.RegisterOutput] {
	return &registerProfessorUsecase{
		registerService:     registerService,
		professorRepository: professorRepository,
	}
}

// Execute implements port.Usecase
func (r *registerProfessorUsecase) Execute(input *inputoutput.RegisterProfessorInput) (*inputoutput.RegisterOutput, error) {
	ctx := context.Background()
	email := input.Email
	name := input.Name
	password := input.Password
	university := input.University
	department := input.Department

	createdUser, err := r.registerService.Execute(ctx, &inputoutput.RegisterInput{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     string(constants.Professor),
	})
	if err != nil {
		return nil, err
	}

	// Save the professor data to database
	newProfessor := &entities.Professor{
		User:       createdUser,
		University: university,
		Department: department,
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

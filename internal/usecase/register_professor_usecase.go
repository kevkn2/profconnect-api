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
func (r *registerProfessorUsecase) Execute(ctx context.Context, input *inputoutput.RegisterProfessorInput) (*inputoutput.RegisterOutput, error) {
	createdUser, err := r.registerService.Execute(ctx, &inputoutput.RegisterInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     string(constants.Professor),
	})
	if err != nil {
		return nil, err
	}

	newProfessor := &entities.Professor{
		User:       createdUser,
		University: input.University,
		Department: input.Department,
	}

	if _, err := r.professorRepository.Create(ctx, newProfessor); err != nil {
		return nil, err
	}

	return &inputoutput.RegisterOutput{
		Message: "User registered successfully",
		UserID:  createdUser.ID,
	}, nil
}

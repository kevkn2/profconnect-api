package auth_usecase

import (
	"context"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type registerStudentUsecase struct {
	registerService    port.Service[inputoutput.RegisterInput, entities.User]
	studentsRepository port.StudentsRepository
}

func NewRegisterStudentUsecase(
	registerService port.Service[inputoutput.RegisterInput, entities.User],
	studentsRepository port.StudentsRepository,
) port.Usecase[inputoutput.RegisterStudentInput, inputoutput.RegisterOutput] {
	return &registerStudentUsecase{
		registerService:    registerService,
		studentsRepository: studentsRepository,
	}
}

func (r *registerStudentUsecase) Execute(ctx context.Context, input *inputoutput.RegisterStudentInput) (*inputoutput.RegisterOutput, error) {
	createdUser, err := r.registerService.Execute(ctx, &inputoutput.RegisterInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     string(constants.Student),
	})
	if err != nil {
		return nil, err
	}

	newStudent := &entities.Student{
		User:              createdUser,
		University:        input.University,
		Department:        input.Department,
		ResearchInterests: input.ResearchInterests,
	}

	if _, err := r.studentsRepository.CreateStudent(ctx, newStudent); err != nil {
		return nil, err
	}

	return &inputoutput.RegisterOutput{
		Message: "User registered successfully",
		UserID:  createdUser.ID,
	}, nil
}

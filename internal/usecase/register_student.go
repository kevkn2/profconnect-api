package usecase

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

// NewRegisterStudentUsecase creates a new instance of RegisterStudentUsecase
func NewRegisterStudentUsecase(
	registerService port.Service[inputoutput.RegisterInput, entities.User],
	studentsRepository port.StudentsRepository,
) port.Usecase[inputoutput.RegisterStudentInput, inputoutput.RegisterOutput] {
	return &registerStudentUsecase{
		registerService:    registerService,
		studentsRepository: studentsRepository,
	}
}

// Execute implements port.Usecase
func (r *registerStudentUsecase) Execute(input *inputoutput.RegisterStudentInput) (*inputoutput.RegisterOutput, error) {
	ctx := context.Background()
	email := input.Email
	name := input.Name
	password := input.Password
	university := input.University
	department := input.Department
	researchInterests := input.ResearchInterests

	createdUser, err := r.registerService.Execute(ctx, &inputoutput.RegisterInput{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     string(constants.Student),
	})
	if err != nil {
		return nil, err
	}

	// save the student data to database
	newStudent := &entities.Student{
		User:              createdUser,
		University:        university,
		Department:        department,
		ResearchInterests: researchInterests,
	}

	// Save student to database
	_, err = r.studentsRepository.CreateStudent(ctx, newStudent)
	if err != nil {
		return nil, err
	}

	return &inputoutput.RegisterOutput{
		Message: "User registered successfully",
		UserID:  createdUser.ID,
	}, nil
}

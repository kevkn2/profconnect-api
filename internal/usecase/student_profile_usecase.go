package usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type studentProfileUsecase struct {
	studentRepository port.StudentsRepository
}

// NewStudentProfileUsecase creates a new instance of StudentProfileUsecase.
func NewStudentProfileUsecase(studentRepository port.StudentsRepository) port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput] {
	return &studentProfileUsecase{
		studentRepository: studentRepository,
	}
}

// Execute implements port.Usecase.
func (u *studentProfileUsecase) Execute(input *inputoutput.ProfileInput) (*inputoutput.StudentProfileOutput, error) {
	ctx := context.Background()

	student, err := u.studentRepository.GetStudentByUserID(ctx, input.UserID)
	if err != nil {
		return nil, domain.InternalErr("failed to retrieve student profile", err)
	}
	if student == nil || student.User == nil {
		return nil, domain.NotFound("student profile not found")
	}

	role := student.User.Role
	if role == "" {
		role = constants.Student
	}

	return &inputoutput.StudentProfileOutput{
		ID:                student.ID,
		UserID:            student.User.ID,
		Name:              student.User.Name,
		Email:             student.User.Email,
		Role:              string(role),
		University:        student.University,
		Department:        student.Department,
		ResearchInterests: student.ResearchInterests,
	}, nil
}

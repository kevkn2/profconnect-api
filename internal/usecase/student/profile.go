package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type profileUsecase struct {
	studentRepository port.StudentsRepository
}

func NewProfileUsecase(studentRepository port.StudentsRepository) port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput] {
	return &profileUsecase{
		studentRepository: studentRepository,
	}
}

func (u *profileUsecase) Execute(ctx context.Context, input *inputoutput.ProfileInput) (*inputoutput.StudentProfileOutput, error) {
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

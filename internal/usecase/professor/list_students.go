package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type listStudentsUsecase struct {
	studentRepository port.StudentsRepository
}

func NewListStudentsUsecase(
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.ListStudentsInput, inputoutput.ListStudentsOutput] {
	return &listStudentsUsecase{
		studentRepository: studentRepository,
	}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, _ *inputoutput.ListStudentsInput) (*inputoutput.ListStudentsOutput, error) {
	students, err := u.studentRepository.ListStudents(ctx)
	if err != nil {
		return nil, domain.InternalErr("failed to list students", err)
	}

	out := make([]*inputoutput.ProjectStudentBrief, 0, len(students))
	for _, s := range students {
		brief := &inputoutput.ProjectStudentBrief{
			StudentID:         s.ID,
			University:        s.University,
			Department:        s.Department,
			ResearchInterests: s.ResearchInterests,
		}
		if s.User != nil {
			brief.UserID = s.User.ID
			brief.Name = s.User.Name
			brief.Email = s.User.Email
		}
		out = append(out, brief)
	}
	return &inputoutput.ListStudentsOutput{Students: out}, nil
}

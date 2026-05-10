package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type profileUsecase struct {
	professorRepository port.ProfessorRepository
}

func NewProfileUsecase(professorRepository port.ProfessorRepository) port.Usecase[inputoutput.ProfileInput, inputoutput.ProfessorProfileOutput] {
	return &profileUsecase{
		professorRepository: professorRepository,
	}
}

func (u *profileUsecase) Execute(ctx context.Context, input *inputoutput.ProfileInput) (*inputoutput.ProfessorProfileOutput, error) {
	professor, err := u.professorRepository.GetByUserID(ctx, input.UserID)
	if err != nil {
		return nil, domain.InternalErr("failed to retrieve professor profile", err)
	}
	if professor == nil || professor.User == nil {
		return nil, domain.NotFound("professor profile not found")
	}

	return &inputoutput.ProfessorProfileOutput{
		ID:         professor.ID,
		UserID:     professor.User.ID,
		Name:       professor.User.Name,
		Email:      professor.User.Email,
		Role:       string(professor.User.Role),
		University: professor.University,
		Department: professor.Department,
	}, nil
}

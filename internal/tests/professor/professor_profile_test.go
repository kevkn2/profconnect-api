package professor

import (
	"context"
	"errors"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	professor_usecase "profconnect-api/internal/usecase/professor"
)

func TestProfessorProfile_Success(t *testing.T) {
	professor := &entities.Professor{
		ID:         "prof-1",
		User:       &entities.User{ID: "user-1", Name: "Prof", Email: "p@x", Role: constants.Professor},
		University: "MIT",
		Department: "CS",
	}
	profRepo := &mocks.ProfessorRepository{}
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }

	uc := professor_usecase.NewProfileUsecase(profRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ProfileInput{UserID: professor.User.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != professor.ID || out.University != "MIT" || out.Department != "CS" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestProfessorProfile_NotFound(t *testing.T) {
	profRepo := &mocks.ProfessorRepository{}
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return nil, nil }

	uc := professor_usecase.NewProfileUsecase(profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ProfileInput{UserID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

func TestProfessorProfile_RepoError(t *testing.T) {
	profRepo := &mocks.ProfessorRepository{}
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
		return nil, errors.New("db down")
	}

	uc := professor_usecase.NewProfileUsecase(profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ProfileInput{UserID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

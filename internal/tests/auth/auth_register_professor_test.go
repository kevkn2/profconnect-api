package auth

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
	auth_usecase "profconnect-api/internal/usecase/auth"
)

func TestRegisterProfessor_Success(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	profRepo := &mocks.ProfessorRepository{}

	createdUser := &entities.User{ID: "user-1", Email: "p@x", Role: constants.Professor}
	core.ExecuteFn = func(_ context.Context, in *inputoutput.RegisterInput) (*entities.User, error) {
		if in.Role != string(constants.Professor) {
			t.Fatalf("expected professor role override, got %q", in.Role)
		}
		return createdUser, nil
	}
	created := false
	profRepo.CreateFn = func(_ context.Context, p *entities.Professor) (*entities.Professor, error) {
		created = true
		if p.User != createdUser {
			t.Fatalf("expected user to be linked to professor record")
		}
		if p.University != "MIT" || p.Department != "CS" {
			t.Fatalf("unexpected professor fields: %+v", p)
		}
		p.ID = "prof-1"
		return p, nil
	}

	uc := auth_usecase.NewRegisterProfessorUsecase(core, profRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.RegisterProfessorInput{
		RegisterInput: inputoutput.RegisterInput{Name: "P", Email: "p@x", Password: "x"},
		University:    "MIT",
		Department:    "CS",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatalf("expected professor record to be created")
	}
	if out.UserID != createdUser.ID {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestRegisterProfessor_CoreError(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	core.ExecuteFn = func(_ context.Context, _ *inputoutput.RegisterInput) (*entities.User, error) {
		return nil, domain.Conflict("email taken")
	}

	uc := auth_usecase.NewRegisterProfessorUsecase(core, &mocks.ProfessorRepository{})
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterProfessorInput{
		RegisterInput: inputoutput.RegisterInput{Email: "p@x"},
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestRegisterProfessor_RepoError(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	profRepo := &mocks.ProfessorRepository{}

	core.ExecuteFn = func(_ context.Context, _ *inputoutput.RegisterInput) (*entities.User, error) {
		return &entities.User{ID: "u"}, nil
	}
	profRepo.CreateFn = func(_ context.Context, _ *entities.Professor) (*entities.Professor, error) {
		return nil, errors.New("db down")
	}

	uc := auth_usecase.NewRegisterProfessorUsecase(core, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterProfessorInput{
		RegisterInput: inputoutput.RegisterInput{Email: "p@x"},
	})
	if err == nil {
		t.Fatalf("expected error from repository to propagate")
	}
}

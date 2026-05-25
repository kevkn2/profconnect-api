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

func TestRegisterStudent_Success(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	stuRepo := &mocks.StudentsRepository{}

	createdUser := &entities.User{ID: "user-1", Email: "s@x", Role: constants.Student}
	core.ExecuteFn = func(_ context.Context, in *inputoutput.RegisterInput) (*entities.User, error) {
		if in.Role != string(constants.Student) {
			t.Fatalf("expected student role override, got %q", in.Role)
		}
		return createdUser, nil
	}
	created := false
	stuRepo.CreateStudentFn = func(_ context.Context, s *entities.Student) (*entities.Student, error) {
		created = true
		if s.User != createdUser {
			t.Fatalf("expected user to be linked to student record")
		}
		if s.University != "MIT" || s.Department != "CS" || s.ResearchInterests != "AI" {
			t.Fatalf("unexpected student fields: %+v", s)
		}
		s.ID = "stu-1"
		return s, nil
	}

	uc := auth_usecase.NewRegisterStudentUsecase(core, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.RegisterStudentInput{
		RegisterInput:     inputoutput.RegisterInput{Name: "S", Email: "s@x", Password: "x"},
		University:        "MIT",
		Department:        "CS",
		ResearchInterests: "AI",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatalf("expected student record to be created")
	}
	if out.UserID != createdUser.ID {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestRegisterStudent_CoreError(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	core.ExecuteFn = func(_ context.Context, _ *inputoutput.RegisterInput) (*entities.User, error) {
		return nil, domain.Conflict("email taken")
	}

	uc := auth_usecase.NewRegisterStudentUsecase(core, &mocks.StudentsRepository{})
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterStudentInput{
		RegisterInput: inputoutput.RegisterInput{Email: "s@x"},
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestRegisterStudent_RepoError(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	stuRepo := &mocks.StudentsRepository{}

	core.ExecuteFn = func(_ context.Context, _ *inputoutput.RegisterInput) (*entities.User, error) {
		return &entities.User{ID: "u"}, nil
	}
	stuRepo.CreateStudentFn = func(_ context.Context, _ *entities.Student) (*entities.Student, error) {
		return nil, errors.New("db down")
	}

	uc := auth_usecase.NewRegisterStudentUsecase(core, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterStudentInput{
		RegisterInput: inputoutput.RegisterInput{Email: "s@x"},
	})
	if err == nil {
		t.Fatalf("expected error from repository to propagate")
	}
}

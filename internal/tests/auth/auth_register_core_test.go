package auth

import (
	"context"
	"errors"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	auth_usecase "profconnect-api/internal/usecase/auth"
)

func TestRegisterCore_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	hasher := &mocks.PasswordHasher{}

	userRepo.GetByEmailFn = func(_ context.Context, _ string) (*entities.User, error) { return nil, nil }
	hasher.HashFn = func(plaintext string) (string, error) {
		if plaintext != "secret" {
			t.Fatalf("unexpected plaintext: %q", plaintext)
		}
		return "hashed", nil
	}
	userRepo.CreateFn = func(_ context.Context, u *entities.User) (*entities.User, error) {
		if u.HashedPassword != "hashed" {
			t.Fatalf("expected hashed password, got %q", u.HashedPassword)
		}
		u.ID = "user-1"
		return u, nil
	}

	uc := auth_usecase.NewRegisterCoreUsecase(userRepo, hasher)
	out, err := uc.Execute(context.Background(), &inputoutput.RegisterInput{
		Name: "Alice", Email: "a@b", Password: "secret", Role: "student",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != "user-1" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestRegisterCore_EmailTaken(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	userRepo.GetByEmailFn = func(_ context.Context, _ string) (*entities.User, error) {
		return &entities.User{ID: "existing"}, nil
	}

	uc := auth_usecase.NewRegisterCoreUsecase(userRepo, &mocks.PasswordHasher{})
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterInput{Email: "taken@x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestRegisterCore_HashError(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	hasher := &mocks.PasswordHasher{}
	userRepo.GetByEmailFn = func(_ context.Context, _ string) (*entities.User, error) { return nil, nil }
	hasher.HashFn = func(_ string) (string, error) { return "", errors.New("boom") }

	uc := auth_usecase.NewRegisterCoreUsecase(userRepo, hasher)
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterInput{Email: "a@b", Password: "p"})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

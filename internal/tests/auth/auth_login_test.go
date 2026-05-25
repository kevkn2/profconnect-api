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

func TestLogin_Success(t *testing.T) {
	user := &entities.User{
		ID:             "user-1",
		Email:          "alice@example.com",
		HashedPassword: "hashed",
		Role:           constants.Student,
	}
	userRepo := &mocks.UserRepository{}
	hasher := &mocks.PasswordHasher{}
	tokens := &mocks.TokenService{}

	userRepo.GetByEmailFn = func(_ context.Context, _ string) (*entities.User, error) { return user, nil }
	hasher.VerifyFn = func(hashed, plaintext string) error {
		if hashed != "hashed" || plaintext != "secret" {
			t.Fatalf("unexpected args to Verify: %q %q", hashed, plaintext)
		}
		return nil
	}
	tokens.GenerateAccessFn = func(uid, email, role string) (string, error) {
		if uid != user.ID || email != user.Email || role != string(user.Role) {
			t.Fatalf("unexpected access token args")
		}
		return "access-tok", nil
	}
	tokens.GenerateRefreshFn = func(_ string) (string, error) { return "refresh-tok", nil }

	uc := auth_usecase.NewLoginUsecase(userRepo, hasher, tokens)
	out, err := uc.Execute(context.Background(), &inputoutput.LoginInput{Email: user.Email, Password: "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.AccessToken != "access-tok" || out.RefreshToken != "refresh-tok" || out.Type != "Bearer" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	userRepo.GetByEmailFn = func(_ context.Context, _ string) (*entities.User, error) { return nil, nil }

	uc := auth_usecase.NewLoginUsecase(userRepo, &mocks.PasswordHasher{}, &mocks.TokenService{})
	_, err := uc.Execute(context.Background(), &inputoutput.LoginInput{Email: "x@x", Password: "y"})
	utils.AssertAppErr(t, err, domain.ErrorTypeUnauthorized)
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	hasher := &mocks.PasswordHasher{}
	userRepo.GetByEmailFn = func(_ context.Context, _ string) (*entities.User, error) {
		return &entities.User{ID: "u", HashedPassword: "h"}, nil
	}
	hasher.VerifyFn = func(_, _ string) error { return errors.New("mismatch") }

	uc := auth_usecase.NewLoginUsecase(userRepo, hasher, &mocks.TokenService{})
	_, err := uc.Execute(context.Background(), &inputoutput.LoginInput{Email: "x@x", Password: "y"})
	utils.AssertAppErr(t, err, domain.ErrorTypeUnauthorized)
}

func TestLogin_RepoError(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	userRepo.GetByEmailFn = func(_ context.Context, _ string) (*entities.User, error) {
		return nil, errors.New("db down")
	}

	uc := auth_usecase.NewLoginUsecase(userRepo, &mocks.PasswordHasher{}, &mocks.TokenService{})
	_, err := uc.Execute(context.Background(), &inputoutput.LoginInput{Email: "x@x", Password: "y"})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

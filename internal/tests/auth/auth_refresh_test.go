package auth

import (
	"context"
	"errors"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	auth_usecase "profconnect-api/internal/usecase/auth"
)

func TestRefresh_Success(t *testing.T) {
	user := &entities.User{ID: "user-1", Email: "a@b", Role: constants.Student}
	userRepo := &mocks.UserRepository{}
	tokens := &mocks.TokenService{}

	tokens.VerifyRefreshFn = func(_ string) (*port.RefreshClaims, error) {
		return &port.RefreshClaims{UserID: user.ID}, nil
	}
	userRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.User, error) { return user, nil }
	tokens.GenerateAccessFn = func(_, _, _ string) (string, error) { return "new-access", nil }
	tokens.GenerateRefreshFn = func(_ string) (string, error) { return "new-refresh", nil }

	uc := auth_usecase.NewRefreshUsecase(userRepo, tokens)
	out, err := uc.Execute(context.Background(), &inputoutput.RefreshInput{RefreshToken: "old"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.AccessToken != "new-access" || out.RefreshToken != "new-refresh" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestRefresh_EmptyToken(t *testing.T) {
	uc := auth_usecase.NewRefreshUsecase(&mocks.UserRepository{}, &mocks.TokenService{})
	_, err := uc.Execute(context.Background(), &inputoutput.RefreshInput{})
	utils.AssertAppErr(t, err, domain.ErrorTypeBadRequest)
}

func TestRefresh_InvalidToken(t *testing.T) {
	tokens := &mocks.TokenService{}
	tokens.VerifyRefreshFn = func(_ string) (*port.RefreshClaims, error) { return nil, errors.New("invalid") }

	uc := auth_usecase.NewRefreshUsecase(&mocks.UserRepository{}, tokens)
	_, err := uc.Execute(context.Background(), &inputoutput.RefreshInput{RefreshToken: "bad"})
	utils.AssertAppErr(t, err, domain.ErrorTypeUnauthorized)
}

func TestRefresh_UserMissing(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	tokens := &mocks.TokenService{}
	tokens.VerifyRefreshFn = func(_ string) (*port.RefreshClaims, error) { return &port.RefreshClaims{UserID: "gone"}, nil }
	userRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.User, error) { return nil, nil }

	uc := auth_usecase.NewRefreshUsecase(userRepo, tokens)
	_, err := uc.Execute(context.Background(), &inputoutput.RefreshInput{RefreshToken: "ok"})
	utils.AssertAppErr(t, err, domain.ErrorTypeUnauthorized)
}

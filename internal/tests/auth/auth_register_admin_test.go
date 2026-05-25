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

func TestRegisterAdmin_Success(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	core.ExecuteFn = func(_ context.Context, in *inputoutput.RegisterInput) (*entities.User, error) {
		if in.Role != string(constants.Admin) {
			t.Fatalf("expected admin role override, got %q", in.Role)
		}
		return &entities.User{ID: "user-1", Email: in.Email}, nil
	}

	uc := auth_usecase.NewRegisterAdminUsecase(core)
	out, err := uc.Execute(context.Background(), &inputoutput.RegisterInput{
		Name: "Adm", Email: "a@x", Password: "p", Role: "ignored",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.UserID != "user-1" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestRegisterAdmin_PropagatesCoreError(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	core.ExecuteFn = func(_ context.Context, _ *inputoutput.RegisterInput) (*entities.User, error) {
		return nil, domain.Conflict("email taken")
	}

	uc := auth_usecase.NewRegisterAdminUsecase(core)
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterInput{Email: "a@x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestRegisterAdmin_PropagatesRawError(t *testing.T) {
	core := &mocks.RegisterCoreService{}
	core.ExecuteFn = func(_ context.Context, _ *inputoutput.RegisterInput) (*entities.User, error) {
		return nil, errors.New("plain")
	}

	uc := auth_usecase.NewRegisterAdminUsecase(core)
	_, err := uc.Execute(context.Background(), &inputoutput.RegisterInput{})
	if err == nil {
		t.Fatalf("expected error to propagate")
	}
}

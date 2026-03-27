package usecase

import (
	"context"
	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	profconnect_utils "profconnect-api/internal/infrastructure/pkg/utils"
)

type LoginUseCase struct {
	userRepository port.UserRepository
}

func NewLoginUseCase(userRepository port.UserRepository) port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput] {
	return &LoginUseCase{
		userRepository: userRepository,
	}
}

// Execute implements port.Usecase.
func (l *LoginUseCase) Execute(input *inputoutput.LoginInput) (*inputoutput.LoginOutput, error) {
	ctx := context.Background()
	email := input.Email
	password := input.Password

	user, err := l.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.InternalErr("failed to retrieve user", err)
	}

	if user == nil {
		return nil, domain.Unauthorized("invalid email or password")
	}

	if err := profconnect_utils.VerifyPassword(user.HashedPassword, password); err != nil {
		return nil, domain.Unauthorized("invalid email or password")
	}

	token, err := profconnect_utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, domain.InternalErr("failed to generate token", err)
	}

	return &inputoutput.LoginOutput{
		Token: token,
		Type:  "Bearer",
	}, nil
}

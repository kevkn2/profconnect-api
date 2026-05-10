package usecase

import (
	"context"
	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	profconnect_utils "profconnect-api/internal/infrastructure/pkg/utils"
)

type LoginUsecase struct {
	userRepository port.UserRepository
}

func NewLoginUsecase(userRepository port.UserRepository) port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput] {
	return &LoginUsecase{
		userRepository: userRepository,
	}
}

// Execute implements port.Usecase.
func (l *LoginUsecase) Execute(ctx context.Context, input *inputoutput.LoginInput) (*inputoutput.LoginOutput, error) {
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

	token, err := profconnect_utils.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, domain.InternalErr("failed to generate token", err)
	}

	refreshToken, err := profconnect_utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to generate refresh token", err)
	}

	return &inputoutput.LoginOutput{
		Token:        token,
		RefreshToken: refreshToken,
		Type:         "Bearer",
	}, nil
}

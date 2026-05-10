package auth_usecase

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

func (l *LoginUsecase) Execute(ctx context.Context, input *inputoutput.LoginInput) (*inputoutput.LoginOutput, error) {
	user, err := l.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.InternalErr("failed to retrieve user", err)
	}

	if user == nil {
		return nil, domain.Unauthorized("invalid email or password")
	}

	if err := profconnect_utils.VerifyPassword(user.HashedPassword, input.Password); err != nil {
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
		AccessToken:  token,
		RefreshToken: refreshToken,
		Role:         string(user.Role),
		Type:         "Bearer",
	}, nil
}

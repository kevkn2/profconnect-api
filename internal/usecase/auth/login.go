package auth_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type LoginUsecase struct {
	userRepository port.UserRepository
	passwordHasher port.PasswordHasher
	tokenService   port.TokenService
}

func NewLoginUsecase(
	userRepository port.UserRepository,
	passwordHasher port.PasswordHasher,
	tokenService port.TokenService,
) port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput] {
	return &LoginUsecase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
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

	if err := l.passwordHasher.Verify(user.HashedPassword, input.Password); err != nil {
		return nil, domain.Unauthorized("invalid email or password")
	}

	accessToken, err := l.tokenService.GenerateAccess(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, domain.InternalErr("failed to generate token", err)
	}

	refreshToken, err := l.tokenService.GenerateRefresh(user.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to generate refresh token", err)
	}

	return &inputoutput.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         string(user.Role),
		Type:         "Bearer",
	}, nil
}

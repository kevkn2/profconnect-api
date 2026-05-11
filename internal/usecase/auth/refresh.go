package auth_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type refreshUsecase struct {
	userRepository port.UserRepository
	tokenService   port.TokenService
}

func NewRefreshUsecase(
	userRepository port.UserRepository,
	tokenService port.TokenService,
) port.Usecase[inputoutput.RefreshInput, inputoutput.RefreshOutput] {
	return &refreshUsecase{
		userRepository: userRepository,
		tokenService:   tokenService,
	}
}

func (r *refreshUsecase) Execute(ctx context.Context, input *inputoutput.RefreshInput) (*inputoutput.RefreshOutput, error) {
	if input.RefreshToken == "" {
		return nil, domain.BadRequest("refresh_token is required")
	}

	claims, err := r.tokenService.VerifyRefresh(input.RefreshToken)
	if err != nil {
		return nil, domain.Unauthorized("invalid or expired refresh token")
	}

	user, err := r.userRepository.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.InternalErr("failed to retrieve user", err)
	}
	if user == nil {
		return nil, domain.Unauthorized("user no longer exists")
	}

	accessToken, err := r.tokenService.GenerateAccess(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, domain.InternalErr("failed to generate access token", err)
	}

	newRefreshToken, err := r.tokenService.GenerateRefresh(user.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to generate refresh token", err)
	}

	return &inputoutput.RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Type:         "Bearer",
	}, nil
}

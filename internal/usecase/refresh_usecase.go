package usecase

import (
	"context"

	"profconnect-api/internal/domain"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	profconnect_utils "profconnect-api/internal/infrastructure/pkg/utils"
)

type refreshUsecase struct {
	userRepository port.UserRepository
}

// NewRefreshUsecase creates a new instance of RefreshUsecase.
func NewRefreshUsecase(userRepository port.UserRepository) port.Usecase[inputoutput.RefreshInput, inputoutput.RefreshOutput] {
	return &refreshUsecase{
		userRepository: userRepository,
	}
}

// Execute implements port.Usecase. It verifies the supplied refresh token,
// reloads the user (so a stale refresh token cannot grant access after a
// role change or deletion), and mints a new access + refresh pair.
func (r *refreshUsecase) Execute(ctx context.Context, input *inputoutput.RefreshInput) (*inputoutput.RefreshOutput, error) {
	if input.RefreshToken == "" {
		return nil, domain.BadRequest("refresh_token is required")
	}

	claims, err := profconnect_utils.VerifyRefreshToken(input.RefreshToken)
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

	accessToken, err := profconnect_utils.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, domain.InternalErr("failed to generate access token", err)
	}

	newRefreshToken, err := profconnect_utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to generate refresh token", err)
	}

	return &inputoutput.RefreshOutput{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		Type:         "Bearer",
	}, nil
}

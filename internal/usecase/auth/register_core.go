package auth_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type registerCoreUsecase struct {
	userRepository port.UserRepository
	passwordHasher port.PasswordHasher
}

func NewRegisterCoreUsecase(
	userRepository port.UserRepository,
	passwordHasher port.PasswordHasher,
) port.Service[inputoutput.RegisterInput, entities.User] {
	return &registerCoreUsecase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (r *registerCoreUsecase) Execute(ctx context.Context, input *inputoutput.RegisterInput) (*entities.User, error) {
	existing, err := r.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.InternalErr("failed to check existing user", err)
	}
	if existing != nil {
		return nil, domain.Conflict("user with this email already exists")
	}

	hashedPassword, err := r.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, domain.InternalErr("failed to hash password", err)
	}

	newUser := &entities.User{
		Name:           input.Name,
		Email:          input.Email,
		HashedPassword: hashedPassword,
		Role:           constants.Roles(input.Role),
	}

	createdUser, err := r.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, domain.InternalErr("failed to create user", err)
	}

	return createdUser, nil
}

package usecase

import (
	"context"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type registerAdminUsecase struct {
	registerService port.Service[inputoutput.RegisterInput, entities.User]
}

// NewRegisterAdminUsecase creates a new instance of RegisterAdminUsecase
func NewRegisterAdminUsecase(
	registerService port.Service[inputoutput.RegisterInput, entities.User],
) port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput] {
	return &registerAdminUsecase{
		registerService: registerService,
	}
}

// Execute implements port.Usecase
func (r *registerAdminUsecase) Execute(ctx context.Context, input *inputoutput.RegisterInput) (*inputoutput.RegisterOutput, error) {
	createdUser, err := r.registerService.Execute(ctx, &inputoutput.RegisterInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     string(constants.Admin),
	})
	if err != nil {
		return nil, err
	}

	return &inputoutput.RegisterOutput{
		Message: "User registered successfully",
		UserID:  createdUser.ID,
	}, nil
}

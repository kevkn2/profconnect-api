package handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler contains all HTTP handlers
type Handler struct {
	registerUsecase port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput]
	loginUseCase    port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput]
}

// NewHandler creates a new handler instance
func NewHandler(
	registerUsecase port.Usecase[inputoutput.RegisterInput, inputoutput.RegisterOutput],
	loginUseCase port.Usecase[inputoutput.LoginInput, inputoutput.LoginOutput],
) *Handler {
	return &Handler{
		registerUsecase: registerUsecase,
		loginUseCase:    loginUseCase,
	}
}

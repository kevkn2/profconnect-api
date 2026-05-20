package student_handler

import (
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

// Handler aggregates HTTP handlers for student-only endpoints.
type Handler struct {
	profileUsecase             port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput]
	applyProjectUsecase        port.Usecase[inputoutput.ApplyProjectInput, inputoutput.ProjectApplicationOutput]
	withdrawApplicationUsecase port.Usecase[inputoutput.WithdrawApplicationInput, inputoutput.EmptyOutput]
	listMyApplicationsUsecase  port.Usecase[inputoutput.ListMyApplicationsInput, inputoutput.ListApplicationsOutput]
}

// New creates a new student handler.
func New(
	profileUsecase port.Usecase[inputoutput.ProfileInput, inputoutput.StudentProfileOutput],
	applyProjectUsecase port.Usecase[inputoutput.ApplyProjectInput, inputoutput.ProjectApplicationOutput],
	withdrawApplicationUsecase port.Usecase[inputoutput.WithdrawApplicationInput, inputoutput.EmptyOutput],
	listMyApplicationsUsecase port.Usecase[inputoutput.ListMyApplicationsInput, inputoutput.ListApplicationsOutput],
) *Handler {
	return &Handler{
		profileUsecase:             profileUsecase,
		applyProjectUsecase:        applyProjectUsecase,
		withdrawApplicationUsecase: withdrawApplicationUsecase,
		listMyApplicationsUsecase:  listMyApplicationsUsecase,
	}
}

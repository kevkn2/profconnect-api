package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type reviewApplicationUsecase struct {
	projectRepository     port.ProjectRepository
	applicationRepository port.ProjectApplicationRepository
	professorRepository   port.ProfessorRepository
}

func NewReviewApplicationUsecase(
	projectRepository port.ProjectRepository,
	applicationRepository port.ProjectApplicationRepository,
	professorRepository port.ProfessorRepository,
) port.Usecase[inputoutput.ReviewApplicationInput, inputoutput.ProjectApplicationOutput] {
	return &reviewApplicationUsecase{
		projectRepository:     projectRepository,
		applicationRepository: applicationRepository,
		professorRepository:   professorRepository,
	}
}

func (u *reviewApplicationUsecase) Execute(ctx context.Context, input *inputoutput.ReviewApplicationInput) (*inputoutput.ProjectApplicationOutput, error) {
	requested := constants.ApplicationStatus(input.Status)
	if requested != constants.ApplicationStatusApproved && requested != constants.ApplicationStatusRejected {
		return nil, domain.ValidationError("status must be 'approved' or 'rejected'")
	}

	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}

	professor, err := u.professorRepository.GetByUserID(ctx, input.ProfessorUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load professor", err)
	}
	if professor == nil {
		return nil, domain.NotFound("professor profile not found")
	}

	if project.Professor == nil || project.Professor.ID != professor.ID {
		return nil, domain.Forbidden("you are not the owner of this project")
	}

	application, err := u.applicationRepository.GetByID(ctx, input.ApplicationID)
	if err != nil {
		return nil, domain.InternalErr("failed to load application", err)
	}
	if application == nil {
		return nil, domain.NotFound("application not found")
	}
	if application.Project.ID != project.ID {
		return nil, domain.NotFound("application does not belong to this project")
	}

	if application.Status != constants.ApplicationStatusPending {
		return nil, domain.Conflict("application has already been reviewed")
	}

	if requested == constants.ApplicationStatusApproved {
		approvedCount, err := u.projectRepository.CountApprovedApplications(ctx, project.ID)
		if err != nil {
			return nil, domain.InternalErr("failed to count approvals", err)
		}
		if approvedCount >= project.Slots {
			return nil, domain.Conflict("all slots have already been filled")
		}
	}

	updated, err := u.applicationRepository.UpdateStatus(ctx, application.ID, requested)
	if err != nil {
		return nil, domain.InternalErr("failed to update application", err)
	}

	if requested == constants.ApplicationStatusApproved {
		approvedCount, err := u.projectRepository.CountApprovedApplications(ctx, project.ID)
		if err != nil {
			return nil, domain.InternalErr("failed to recount approvals", err)
		}
		if approvedCount >= project.Slots {
			if err := u.projectRepository.UpdateStatus(ctx, project.ID, constants.ProjectStatusClosed); err != nil {
				return nil, domain.InternalErr("failed to close project", err)
			}
		}
	}

	// Restore student context for the response.
	updated.Student = application.Student
	return &inputoutput.ProjectApplicationOutput{
		ID:        updated.ID,
		Project:   &inputoutput.ProjectShortForApp{
			Title: project.Title,
			Description: project.Description,
			Status: string(project.Status),
		},
		Student:   &inputoutput.ProjectStudentBrief{
			StudentID: updated.Student.ID,
			UserID: updated.Student.User.ID,
			Email: updated.Student.User.Email,
			Name: updated.Student.User.Name,
			University: updated.Student.University,
			Department: updated.Student.Department,
			ResearchInterests: updated.Student.ResearchInterests,
		},
		Status:    string(updated.Status),
		Message:   updated.Message,
	}, nil
}

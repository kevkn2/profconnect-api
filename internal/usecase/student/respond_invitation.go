package student_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type respondInvitationUsecase struct {
	projectRepository    port.ProjectRepository
	invitationRepository port.ProjectInvitationRepository
	memberRepository     port.ProjectMemberRepository
	studentRepository    port.StudentsRepository
}

func NewRespondInvitationUsecase(
	projectRepository port.ProjectRepository,
	invitationRepository port.ProjectInvitationRepository,
	memberRepository port.ProjectMemberRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.RespondInvitationInput, inputoutput.ProjectInvitationOutput] {
	return &respondInvitationUsecase{
		projectRepository:    projectRepository,
		invitationRepository: invitationRepository,
		memberRepository:     memberRepository,
		studentRepository:    studentRepository,
	}
}

func (u *respondInvitationUsecase) Execute(ctx context.Context, input *inputoutput.RespondInvitationInput) (*inputoutput.ProjectInvitationOutput, error) {
	requested := constants.InvitationStatus(input.Status)
	if requested != constants.InvitationStatusAccepted && requested != constants.InvitationStatusDeclined {
		return nil, domain.ValidationError("status must be 'accepted' or 'declined'")
	}

	student, err := u.studentRepository.GetStudentByUserID(ctx, input.StudentUserID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student profile not found")
	}

	invitation, err := u.invitationRepository.GetByID(ctx, input.InvitationID)
	if err != nil {
		return nil, domain.InternalErr("failed to load invitation", err)
	}
	if invitation == nil {
		return nil, domain.NotFound("invitation not found")
	}
	if invitation.Student == nil || invitation.Student.ID != student.ID {
		return nil, domain.Forbidden("this invitation is not addressed to you")
	}
	if invitation.Status != constants.InvitationStatusPending {
		return nil, domain.Conflict("invitation has already been responded to")
	}

	project, err := u.projectRepository.GetByID(ctx, invitation.Project.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}

	if requested == constants.InvitationStatusAccepted {
		if project.Status != constants.ProjectStatusOpen {
			return nil, domain.Conflict("project is no longer accepting new members")
		}
		activeCount, err := u.memberRepository.CountActive(ctx, project.ID)
		if err != nil {
			return nil, domain.InternalErr("failed to count active members", err)
		}
		if activeCount >= project.Slots {
			return nil, domain.Conflict("all slots have already been filled")
		}
	}

	updated, err := u.invitationRepository.UpdateStatus(ctx, invitation.ID, requested)
	if err != nil {
		return nil, domain.InternalErr("failed to update invitation", err)
	}

	if requested == constants.InvitationStatusAccepted {
		if _, err := u.memberRepository.Create(ctx, &entities.ProjectMember{
			Project:     project,
			Student:     student,
			Source:      constants.MemberSourceInvitation,
			SourceRefID: invitation.ID,
			Status:      constants.MemberStatusActive,
		}); err != nil {
			return nil, domain.InternalErr("failed to add project member", err)
		}

		activeCount, err := u.memberRepository.CountActive(ctx, project.ID)
		if err != nil {
			return nil, domain.InternalErr("failed to recount active members", err)
		}
		if activeCount >= project.Slots && project.Status == constants.ProjectStatusOpen {
			if err := u.projectRepository.UpdateStatus(ctx, project.ID, constants.ProjectStatusClosed); err != nil {
				return nil, domain.InternalErr("failed to close project", err)
			}
		}
	}

	updated.Student = student
	return &inputoutput.ProjectInvitationOutput{
		ID:      updated.ID,
		Status:  string(updated.Status),
		Message: updated.Message,
		Student: &inputoutput.ProjectStudentBrief{
			StudentID:         student.ID,
			UserID:            student.User.ID,
			Name:              student.User.Name,
			Email:             student.User.Email,
			University:        student.University,
			Department:        student.Department,
			ResearchInterests: student.ResearchInterests,
		},
		Project: &inputoutput.ProjectShortForApp{
			Title:       project.Title,
			Description: project.Description,
			Status:      string(project.Status),
		},
	}, nil
}

package professor_usecase

import (
	"context"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

type sendInvitationUsecase struct {
	projectRepository     port.ProjectRepository
	invitationRepository  port.ProjectInvitationRepository
	memberRepository      port.ProjectMemberRepository
	applicationRepository port.ProjectApplicationRepository
	professorRepository   port.ProfessorRepository
	studentRepository     port.StudentsRepository
}

func NewSendInvitationUsecase(
	projectRepository port.ProjectRepository,
	invitationRepository port.ProjectInvitationRepository,
	memberRepository port.ProjectMemberRepository,
	applicationRepository port.ProjectApplicationRepository,
	professorRepository port.ProfessorRepository,
	studentRepository port.StudentsRepository,
) port.Usecase[inputoutput.SendInvitationInput, inputoutput.ProjectInvitationOutput] {
	return &sendInvitationUsecase{
		projectRepository:     projectRepository,
		invitationRepository:  invitationRepository,
		memberRepository:      memberRepository,
		applicationRepository: applicationRepository,
		professorRepository:   professorRepository,
		studentRepository:     studentRepository,
	}
}

func (u *sendInvitationUsecase) Execute(ctx context.Context, input *inputoutput.SendInvitationInput) (*inputoutput.ProjectInvitationOutput, error) {
	if input.StudentID == "" {
		return nil, domain.ValidationError("student_id is required")
	}

	project, err := u.projectRepository.GetByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.InternalErr("failed to load project", err)
	}
	if project == nil {
		return nil, domain.NotFound("project not found")
	}
	if project.Status != constants.ProjectStatusOpen {
		return nil, domain.Conflict("project is not open for invitations")
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

	student, err := u.studentRepository.GetStudentByID(ctx, input.StudentID)
	if err != nil {
		return nil, domain.InternalErr("failed to load student", err)
	}
	if student == nil {
		return nil, domain.NotFound("student not found")
	}

	member, err := u.memberRepository.GetActiveByProjectAndStudent(ctx, project.ID, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to check membership", err)
	}
	if member != nil {
		return nil, domain.Conflict("student is already a member of this project")
	}

	existingApp, err := u.applicationRepository.GetByProjectAndStudent(ctx, project.ID, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to check existing application", err)
	}
	if existingApp != nil && (existingApp.Status == constants.ApplicationStatusPending || existingApp.Status == constants.ApplicationStatusApproved) {
		return nil, domain.Conflict("student already has an active application; review it instead of inviting")
	}

	existingInv, err := u.invitationRepository.GetByProjectAndStudent(ctx, project.ID, student.ID)
	if err != nil {
		return nil, domain.InternalErr("failed to check existing invitation", err)
	}
	if existingInv != nil && existingInv.Status == constants.InvitationStatusPending {
		return nil, domain.Conflict("an invitation has already been sent to this student")
	}
	if existingInv != nil {
		return nil, domain.Conflict("this student has a prior invitation on this project; cannot resend")
	}

	created, err := u.invitationRepository.Create(ctx, &entities.ProjectInvitation{
		Project: project,
		Student: student,
		Status:  constants.InvitationStatusPending,
		Message: input.Message,
	})
	if err != nil {
		return nil, domain.InternalErr("failed to create invitation", err)
	}

	return invitationToOutput(created, project, student), nil
}

func invitationToOutput(inv *entities.ProjectInvitation, project *entities.Project, student *entities.Student) *inputoutput.ProjectInvitationOutput {
	out := &inputoutput.ProjectInvitationOutput{
		ID:      inv.ID,
		Status:  string(inv.Status),
		Message: inv.Message,
	}
	if inv.RespondedAt != nil {
		out.RespondedAt = inv.RespondedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	if student != nil {
		brief := &inputoutput.ProjectStudentBrief{
			StudentID:         student.ID,
			University:        student.University,
			Department:        student.Department,
			ResearchInterests: student.ResearchInterests,
		}
		if student.User != nil {
			brief.UserID = student.User.ID
			brief.Name = student.User.Name
			brief.Email = student.User.Email
		}
		out.Student = brief
	}
	if project != nil {
		out.Project = &inputoutput.ProjectShortForApp{
			Title:       project.Title,
			Description: project.Description,
			Status:      string(project.Status),
		}
	}
	return out
}

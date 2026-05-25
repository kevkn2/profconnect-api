package professor

import (
	"context"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	professor_usecase "profconnect-api/internal/usecase/professor"
)

func cancelInvitationFixtures() (project *entities.Project, professor *entities.Professor, invitation *entities.ProjectInvitation) {
	professor = &entities.Professor{ID: "prof-1", User: &entities.User{ID: "user-prof-1"}}
	project = &entities.Project{ID: "proj-1", Professor: professor, Status: constants.ProjectStatusOpen}
	invitation = &entities.ProjectInvitation{
		ID:      "inv-1",
		Project: &entities.Project{ID: project.ID},
		Status:  constants.InvitationStatusPending,
	}
	return
}

func TestCancelInvitation_Success(t *testing.T) {
	project, professor, invitation := cancelInvitationFixtures()
	projectRepo := &mocks.ProjectRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	invRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return invitation, nil }
	cancelled := false
	invRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.InvitationStatus) (*entities.ProjectInvitation, error) {
		cancelled = true
		if s != constants.InvitationStatusCancelled {
			t.Fatalf("expected cancelled, got %q", s)
		}
		return invitation, nil
	}

	uc := professor_usecase.NewCancelInvitationUsecase(projectRepo, invRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CancelInvitationInput{
		ProfessorUserID: professor.User.ID, ProjectID: project.ID, InvitationID: invitation.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cancelled {
		t.Fatalf("expected UpdateStatus to be called")
	}
}

func TestCancelInvitation_NotPending(t *testing.T) {
	project, professor, invitation := cancelInvitationFixtures()
	invitation.Status = constants.InvitationStatusAccepted
	projectRepo := &mocks.ProjectRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	invRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return invitation, nil }

	uc := professor_usecase.NewCancelInvitationUsecase(projectRepo, invRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CancelInvitationInput{
		ProfessorUserID: professor.User.ID, ProjectID: project.ID, InvitationID: invitation.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestCancelInvitation_NotOwner(t *testing.T) {
	project, professor, invitation := cancelInvitationFixtures()
	projectRepo := &mocks.ProjectRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
		return &entities.Professor{ID: "prof-other"}, nil
	}

	uc := professor_usecase.NewCancelInvitationUsecase(projectRepo, &mocks.ProjectInvitationRepository{}, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CancelInvitationInput{
		ProfessorUserID: professor.User.ID, ProjectID: project.ID, InvitationID: invitation.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeForbidden)
}

func TestCancelInvitation_BelongsToDifferentProject(t *testing.T) {
	project, professor, invitation := cancelInvitationFixtures()
	invitation.Project = &entities.Project{ID: "proj-other"}
	projectRepo := &mocks.ProjectRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	invRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return invitation, nil }

	uc := professor_usecase.NewCancelInvitationUsecase(projectRepo, invRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CancelInvitationInput{
		ProfessorUserID: professor.User.ID, ProjectID: project.ID, InvitationID: invitation.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

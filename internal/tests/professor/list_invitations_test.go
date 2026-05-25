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

func TestListInvitationsByProject_Success(t *testing.T) {
	professor := &entities.Professor{ID: "prof-1", User: &entities.User{ID: "user-prof-1"}}
	project := &entities.Project{ID: "proj-1", Professor: professor, Status: constants.ProjectStatusOpen}

	projectRepo := &mocks.ProjectRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	invRepo.ListByProjectFn = func(_ context.Context, _ string) ([]*entities.ProjectInvitation, error) {
		return []*entities.ProjectInvitation{
			{ID: "inv-1", Status: constants.InvitationStatusPending, Student: &entities.Student{ID: "s-1", User: &entities.User{ID: "u-1"}}},
			{ID: "inv-2", Status: constants.InvitationStatusAccepted, Student: &entities.Student{ID: "s-2", User: &entities.User{ID: "u-2"}}},
		}, nil
	}

	uc := professor_usecase.NewListInvitationsByProjectUsecase(projectRepo, invRepo, profRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListInvitationsByProjectInput{
		ProfessorUserID: professor.User.ID, ProjectID: project.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Invitations) != 2 {
		t.Fatalf("expected 2 invitations, got %d", len(out.Invitations))
	}
}

func TestListInvitationsByProject_NotOwner(t *testing.T) {
	project := &entities.Project{ID: "proj-1", Professor: &entities.Professor{ID: "prof-1"}}
	projectRepo := &mocks.ProjectRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
		return &entities.Professor{ID: "prof-other"}, nil
	}

	uc := professor_usecase.NewListInvitationsByProjectUsecase(projectRepo, &mocks.ProjectInvitationRepository{}, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ListInvitationsByProjectInput{ProjectID: project.ID})
	utils.AssertAppErr(t, err, domain.ErrorTypeForbidden)
}

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

func TestListApplicationsByProject_BucketsByStatus(t *testing.T) {
	professor := &entities.Professor{ID: "prof-1", User: &entities.User{ID: "user-prof-1"}}
	project := &entities.Project{ID: "proj-1", Professor: professor, Status: constants.ProjectStatusOpen}
	makeApp := func(id string, status constants.ApplicationStatus) *entities.ProjectApplication {
		return &entities.ProjectApplication{
			ID:      id,
			Project: project,
			Student: &entities.Student{ID: "s-" + id, User: &entities.User{ID: "u-" + id, Name: id, Email: id + "@x"}},
			Status:  status,
		}
	}

	projectRepo := &mocks.ProjectRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	appRepo.ListByProjectFn = func(_ context.Context, _ string) ([]*entities.ProjectApplication, error) {
		return []*entities.ProjectApplication{
			makeApp("1", constants.ApplicationStatusApproved),
			makeApp("2", constants.ApplicationStatusPending),
			makeApp("3", constants.ApplicationStatusRejected),
			makeApp("4", constants.ApplicationStatusApproved),
		}, nil
	}

	uc := professor_usecase.NewListApplicationsByProjectUsecase(projectRepo, appRepo, profRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListApplicationsByProjectInput{
		ProfessorUserID: professor.User.ID, ProjectID: project.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.ApprovedApplications) != 2 {
		t.Fatalf("expected 2 approved, got %d", len(out.ApprovedApplications))
	}
	if len(out.PendingApplications) != 1 {
		t.Fatalf("expected 1 pending, got %d", len(out.PendingApplications))
	}
}

func TestListApplicationsByProject_NotOwner(t *testing.T) {
	project := &entities.Project{ID: "proj-1", Professor: &entities.Professor{ID: "prof-1"}}
	projectRepo := &mocks.ProjectRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
		return &entities.Professor{ID: "prof-other"}, nil
	}

	uc := professor_usecase.NewListApplicationsByProjectUsecase(projectRepo, &mocks.ProjectApplicationRepository{}, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ListApplicationsByProjectInput{
		ProfessorUserID: "u", ProjectID: project.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeForbidden)
}

func TestListApplicationsByProject_ProjectNotFound(t *testing.T) {
	projectRepo := &mocks.ProjectRepository{}
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return nil, nil }

	uc := professor_usecase.NewListApplicationsByProjectUsecase(projectRepo, &mocks.ProjectApplicationRepository{}, &mocks.ProfessorRepository{})
	_, err := uc.Execute(context.Background(), &inputoutput.ListApplicationsByProjectInput{ProjectID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

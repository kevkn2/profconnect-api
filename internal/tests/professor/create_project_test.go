package professor

import (
	"context"
	"errors"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	professor_usecase "profconnect-api/internal/usecase/professor"
)

func TestCreateProject_Success(t *testing.T) {
	professor := &entities.Professor{ID: "prof-1", User: &entities.User{ID: "user-prof-1"}}
	projectRepo := &mocks.ProjectRepository{}
	profRepo := &mocks.ProfessorRepository{}

	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	projectRepo.CreateFn = func(_ context.Context, p *entities.Project) (*entities.Project, error) {
		if p.Status != constants.ProjectStatusOpen {
			t.Fatalf("expected new project to be open, got %q", p.Status)
		}
		p.ID = "proj-1"
		return p, nil
	}

	uc := professor_usecase.NewCreateProjectUsecase(projectRepo, profRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.CreateProjectInput{
		ProfessorUserID: professor.User.ID,
		Title:           "AI Research",
		Description:     "Study LLMs",
		Slots:           3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != "proj-1" || out.Slots != 3 || out.Status != string(constants.ProjectStatusOpen) {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestCreateProject_Validation(t *testing.T) {
	cases := []struct {
		name  string
		input *inputoutput.CreateProjectInput
	}{
		{name: "missing title", input: &inputoutput.CreateProjectInput{Description: "d", Slots: 1}},
		{name: "missing description", input: &inputoutput.CreateProjectInput{Title: "t", Slots: 1}},
		{name: "zero slots", input: &inputoutput.CreateProjectInput{Title: "t", Description: "d", Slots: 0}},
		{name: "negative slots", input: &inputoutput.CreateProjectInput{Title: "t", Description: "d", Slots: -1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc := professor_usecase.NewCreateProjectUsecase(&mocks.ProjectRepository{}, &mocks.ProfessorRepository{})
			_, err := uc.Execute(context.Background(), tc.input)
			utils.AssertAppErr(t, err, domain.ErrorTypeValidation)
		})
	}
}

func TestCreateProject_ProfessorNotFound(t *testing.T) {
	profRepo := &mocks.ProfessorRepository{}
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return nil, nil }

	uc := professor_usecase.NewCreateProjectUsecase(&mocks.ProjectRepository{}, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CreateProjectInput{Title: "t", Description: "d", Slots: 1})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

func TestCreateProject_RepoError(t *testing.T) {
	projectRepo := &mocks.ProjectRepository{}
	profRepo := &mocks.ProfessorRepository{}
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
		return &entities.Professor{ID: "p"}, nil
	}
	projectRepo.CreateFn = func(_ context.Context, _ *entities.Project) (*entities.Project, error) {
		return nil, errors.New("db down")
	}

	uc := professor_usecase.NewCreateProjectUsecase(projectRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CreateProjectInput{Title: "t", Description: "d", Slots: 1})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

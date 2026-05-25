package project

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
	project_usecase "profconnect-api/internal/usecase/project"
)

func TestGetProject_Success(t *testing.T) {
	project := &entities.Project{
		ID:          "proj-1",
		Title:       "AI",
		Description: "d",
		Slots:       3,
		Status:      constants.ProjectStatusOpen,
		Professor: &entities.Professor{
			ID:         "prof-1",
			User:       &entities.User{ID: "user-1", Name: "P", Email: "p@x"},
			University: "MIT", Department: "CS",
		},
	}
	repo := &mocks.ProjectRepository{}
	repo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }

	uc := project_usecase.NewGetProjectUsecase(repo)
	out, err := uc.Execute(context.Background(), &inputoutput.GetProjectInput{ProjectID: project.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != project.ID || out.ProfessorID != project.Professor.ID || out.Professor.Name != "P" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestGetProject_NotFound(t *testing.T) {
	repo := &mocks.ProjectRepository{}
	repo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return nil, nil }

	uc := project_usecase.NewGetProjectUsecase(repo)
	_, err := uc.Execute(context.Background(), &inputoutput.GetProjectInput{ProjectID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

func TestGetProject_RepoError(t *testing.T) {
	repo := &mocks.ProjectRepository{}
	repo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) {
		return nil, errors.New("db down")
	}

	uc := project_usecase.NewGetProjectUsecase(repo)
	_, err := uc.Execute(context.Background(), &inputoutput.GetProjectInput{ProjectID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

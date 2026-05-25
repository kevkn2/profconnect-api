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

func TestListProjects_Success(t *testing.T) {
	makeProject := func(id, title string) *entities.Project {
		return &entities.Project{
			ID: id, Title: title, Description: "d", Slots: 2, Status: constants.ProjectStatusOpen,
			Professor: &entities.Professor{
				ID:   "prof-" + id,
				User: &entities.User{ID: "u-" + id, Name: "P-" + id, Email: id + "@x"},
			},
		}
	}
	repo := &mocks.ProjectRepository{}
	repo.ListFn = func(_ context.Context) ([]*entities.Project, error) {
		return []*entities.Project{makeProject("1", "A"), makeProject("2", "B")}, nil
	}

	uc := project_usecase.NewListProjectsUsecase(repo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListProjectsInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(out.Projects))
	}
}

func TestListProjects_RepoError(t *testing.T) {
	repo := &mocks.ProjectRepository{}
	repo.ListFn = func(_ context.Context) ([]*entities.Project, error) { return nil, errors.New("db down") }

	uc := project_usecase.NewListProjectsUsecase(repo)
	_, err := uc.Execute(context.Background(), &inputoutput.ListProjectsInput{})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

func TestListProjects_Empty(t *testing.T) {
	repo := &mocks.ProjectRepository{}
	repo.ListFn = func(_ context.Context) ([]*entities.Project, error) { return nil, nil }

	uc := project_usecase.NewListProjectsUsecase(repo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListProjectsInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Projects) != 0 {
		t.Fatalf("expected 0 projects, got %d", len(out.Projects))
	}
}

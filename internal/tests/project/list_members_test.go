package project

import (
	"context"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	project_usecase "profconnect-api/internal/usecase/project"
)

func TestListMembersByProject_Success(t *testing.T) {
	project := &entities.Project{ID: "proj-1", Status: constants.ProjectStatusOpen}
	projectRepo := &mocks.ProjectRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	memberRepo.ListByProjectFn = func(_ context.Context, _ string) ([]*entities.ProjectMember, error) {
		return []*entities.ProjectMember{
			{
				ID: "m-1", Status: constants.MemberStatusActive, Source: constants.MemberSourceApplication,
				Student: &entities.Student{ID: "s-1", User: &entities.User{ID: "u-1", Name: "A"}},
			},
			{
				ID: "m-2", Status: constants.MemberStatusLeft, Source: constants.MemberSourceInvitation,
				Student: &entities.Student{ID: "s-2", User: &entities.User{ID: "u-2", Name: "B"}},
			},
		}, nil
	}

	uc := project_usecase.NewListMembersByProjectUsecase(projectRepo, memberRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListMembersByProjectInput{ProjectID: project.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(out.Members))
	}
	if out.Members[0].Source != string(constants.MemberSourceApplication) {
		t.Fatalf("expected source 'application', got %q", out.Members[0].Source)
	}
}

func TestListMembersByProject_ProjectNotFound(t *testing.T) {
	projectRepo := &mocks.ProjectRepository{}
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return nil, nil }

	uc := project_usecase.NewListMembersByProjectUsecase(projectRepo, &mocks.ProjectMemberRepository{})
	_, err := uc.Execute(context.Background(), &inputoutput.ListMembersByProjectInput{ProjectID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

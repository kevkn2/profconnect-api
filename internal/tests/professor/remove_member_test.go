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

func removeMemberFixtures() (project *entities.Project, professor *entities.Professor, member *entities.ProjectMember) {
	professor = &entities.Professor{
		ID:   "prof-1",
		User: &entities.User{ID: "user-prof-1"},
	}
	project = &entities.Project{
		ID:        "proj-1",
		Professor: professor,
		Slots:     2,
		Status:    constants.ProjectStatusOpen,
	}
	member = &entities.ProjectMember{
		ID:      "mem-1",
		Project: &entities.Project{ID: project.ID},
		Student: &entities.Student{ID: "stu-1"},
		Status:  constants.MemberStatusActive,
	}
	return
}

func TestRemoveMember_Success(t *testing.T) {
	project, professor, member := removeMemberFixtures()
	projectRepo := &mocks.ProjectRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	memberRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectMember, error) { return member, nil }
	updateCalled := false
	memberRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.MemberStatus) (*entities.ProjectMember, error) {
		updateCalled = true
		if s != constants.MemberStatusRemoved {
			t.Fatalf("expected status 'removed', got %q", s)
		}
		return member, nil
	}

	uc := professor_usecase.NewRemoveMemberUsecase(projectRepo, memberRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.RemoveMemberInput{
		ProfessorUserID: professor.User.ID,
		ProjectID:       project.ID,
		MemberID:        member.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updateCalled {
		t.Fatalf("expected member UpdateStatus to be called")
	}
}

func TestRemoveMember_ReopensClosedProject(t *testing.T) {
	project, professor, member := removeMemberFixtures()
	project.Status = constants.ProjectStatusClosed
	projectRepo := &mocks.ProjectRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	memberRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectMember, error) { return member, nil }
	memberRepo.UpdateStatusFn = func(_ context.Context, _ string, _ constants.MemberStatus) (*entities.ProjectMember, error) {
		return member, nil
	}
	reopened := false
	projectRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.ProjectStatus) error {
		if s == constants.ProjectStatusOpen {
			reopened = true
		}
		return nil
	}

	uc := professor_usecase.NewRemoveMemberUsecase(projectRepo, memberRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.RemoveMemberInput{
		ProfessorUserID: professor.User.ID,
		ProjectID:       project.ID,
		MemberID:        member.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reopened {
		t.Fatalf("expected closed project to be reopened when slot frees")
	}
}

func TestRemoveMember_ValidationAndPreconditions(t *testing.T) {
	cases := []struct {
		name    string
		setup   func(*mocks.ProjectRepository, *mocks.ProjectMemberRepository, *mocks.ProfessorRepository, *entities.Project, *entities.Professor, *entities.ProjectMember)
		wantErr domain.ErrorType
	}{
		{
			name: "not project owner",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectMemberRepository, profRepo *mocks.ProfessorRepository, p *entities.Project, _ *entities.Professor, _ *entities.ProjectMember) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
					return &entities.Professor{ID: "prof-other"}, nil
				}
			},
			wantErr: domain.ErrorTypeForbidden,
		},
		{
			name: "member belongs to different project",
			setup: func(pr *mocks.ProjectRepository, mr *mocks.ProjectMemberRepository, profRepo *mocks.ProfessorRepository, p *entities.Project, prof *entities.Professor, m *entities.ProjectMember) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return prof, nil }
				m.Project = &entities.Project{ID: "proj-other"}
				mr.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectMember, error) { return m, nil }
			},
			wantErr: domain.ErrorTypeNotFound,
		},
		{
			name: "member already inactive",
			setup: func(pr *mocks.ProjectRepository, mr *mocks.ProjectMemberRepository, profRepo *mocks.ProfessorRepository, p *entities.Project, prof *entities.Professor, m *entities.ProjectMember) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return prof, nil }
				m.Status = constants.MemberStatusLeft
				mr.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectMember, error) { return m, nil }
			},
			wantErr: domain.ErrorTypeConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			project, professor, member := removeMemberFixtures()
			projectRepo := &mocks.ProjectRepository{}
			memberRepo := &mocks.ProjectMemberRepository{}
			profRepo := &mocks.ProfessorRepository{}
			tc.setup(projectRepo, memberRepo, profRepo, project, professor, member)

			uc := professor_usecase.NewRemoveMemberUsecase(projectRepo, memberRepo, profRepo)
			_, err := uc.Execute(context.Background(), &inputoutput.RemoveMemberInput{
				ProfessorUserID: professor.User.ID,
				ProjectID:       project.ID,
				MemberID:        member.ID,
			})
			utils.AssertAppErr(t, err, tc.wantErr)
		})
	}
}

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

func reviewFixtures() (project *entities.Project, professor *entities.Professor, application *entities.ProjectApplication) {
	professor = &entities.Professor{
		ID:   "prof-1",
		User: &entities.User{ID: "user-prof-1", Name: "Prof", Email: "prof@example.com"},
	}
	project = &entities.Project{
		ID:        "proj-1",
		Professor: professor,
		Title:     "AI Research",
		Slots:     2,
		Status:    constants.ProjectStatusOpen,
	}
	student := &entities.Student{
		ID:   "stu-1",
		User: &entities.User{ID: "user-stu-1", Name: "Stu", Email: "stu@example.com"},
	}
	application = &entities.ProjectApplication{
		ID:      "app-1",
		Project: &entities.Project{ID: project.ID},
		Student: student,
		Status:  constants.ApplicationStatusPending,
	}
	return
}

func TestReviewApplication_ApproveCreatesMember(t *testing.T) {
	project, professor, app := reviewFixtures()
	projectRepo := &mocks.ProjectRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }
	memberRepo.CountActiveFn = func(_ context.Context, _ string) (int, error) { return 0, nil }
	appRepo.UpdateStatusFn = func(_ context.Context, _ string, _ constants.ApplicationStatus) (*entities.ProjectApplication, error) {
		return &entities.ProjectApplication{ID: app.ID, Project: app.Project, Student: app.Student, Status: constants.ApplicationStatusApproved}, nil
	}
	memberCreated := false
	memberRepo.CreateFn = func(_ context.Context, m *entities.ProjectMember) (*entities.ProjectMember, error) {
		memberCreated = true
		if m.Source != constants.MemberSourceApplication {
			t.Fatalf("expected source 'application', got %q", m.Source)
		}
		if m.SourceRefID != app.ID {
			t.Fatalf("expected source ref %q, got %q", app.ID, m.SourceRefID)
		}
		return m, nil
	}

	uc := professor_usecase.NewReviewApplicationUsecase(projectRepo, appRepo, memberRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ReviewApplicationInput{
		ProfessorUserID: professor.User.ID,
		ProjectID:       project.ID,
		ApplicationID:   app.ID,
		Status:          string(constants.ApplicationStatusApproved),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !memberCreated {
		t.Fatalf("expected member Create to be called on approval")
	}
}

func TestReviewApplication_ApproveClosesProjectWhenSlotsFill(t *testing.T) {
	project, professor, app := reviewFixtures()
	project.Slots = 1
	projectRepo := &mocks.ProjectRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }
	calls := 0
	memberRepo.CountActiveFn = func(_ context.Context, _ string) (int, error) {
		calls++
		if calls == 1 {
			return 0, nil
		}
		return 1, nil
	}
	appRepo.UpdateStatusFn = func(_ context.Context, _ string, _ constants.ApplicationStatus) (*entities.ProjectApplication, error) {
		return &entities.ProjectApplication{ID: app.ID, Project: app.Project, Student: app.Student, Status: constants.ApplicationStatusApproved}, nil
	}
	memberRepo.CreateFn = func(_ context.Context, m *entities.ProjectMember) (*entities.ProjectMember, error) { return m, nil }
	closed := false
	projectRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.ProjectStatus) error {
		if s == constants.ProjectStatusClosed {
			closed = true
		}
		return nil
	}

	uc := professor_usecase.NewReviewApplicationUsecase(projectRepo, appRepo, memberRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ReviewApplicationInput{
		ProfessorUserID: professor.User.ID,
		ProjectID:       project.ID,
		ApplicationID:   app.ID,
		Status:          string(constants.ApplicationStatusApproved),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !closed {
		t.Fatalf("expected project to be closed when slots fill")
	}
}

func TestReviewApplication_RejectDoesNotCreateMember(t *testing.T) {
	project, professor, app := reviewFixtures()
	projectRepo := &mocks.ProjectRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	profRepo := &mocks.ProfessorRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }
	appRepo.UpdateStatusFn = func(_ context.Context, _ string, _ constants.ApplicationStatus) (*entities.ProjectApplication, error) {
		return &entities.ProjectApplication{ID: app.ID, Project: app.Project, Student: app.Student, Status: constants.ApplicationStatusRejected}, nil
	}

	uc := professor_usecase.NewReviewApplicationUsecase(projectRepo, appRepo, memberRepo, profRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ReviewApplicationInput{
		ProfessorUserID: professor.User.ID,
		ProjectID:       project.ID,
		ApplicationID:   app.ID,
		Status:          string(constants.ApplicationStatusRejected),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReviewApplication_ValidationAndPreconditions(t *testing.T) {
	cases := []struct {
		name    string
		setup   func(*mocks.ProjectRepository, *mocks.ProjectApplicationRepository, *mocks.ProjectMemberRepository, *mocks.ProfessorRepository, *entities.Project, *entities.Professor, *entities.ProjectApplication)
		input   func(*entities.Project, *entities.Professor, *entities.ProjectApplication) *inputoutput.ReviewApplicationInput
		wantErr domain.ErrorType
	}{
		{
			name: "invalid status",
			setup: func(_ *mocks.ProjectRepository, _ *mocks.ProjectApplicationRepository, _ *mocks.ProjectMemberRepository, _ *mocks.ProfessorRepository, _ *entities.Project, _ *entities.Professor, _ *entities.ProjectApplication) {
			},
			input: func(p *entities.Project, prof *entities.Professor, app *entities.ProjectApplication) *inputoutput.ReviewApplicationInput {
				return &inputoutput.ReviewApplicationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, ApplicationID: app.ID, Status: "withdrawn"}
			},
			wantErr: domain.ErrorTypeValidation,
		},
		{
			name: "not project owner",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectApplicationRepository, _ *mocks.ProjectMemberRepository, profRepo *mocks.ProfessorRepository, p *entities.Project, _ *entities.Professor, _ *entities.ProjectApplication) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
					return &entities.Professor{ID: "prof-other"}, nil
				}
			},
			input: func(p *entities.Project, prof *entities.Professor, app *entities.ProjectApplication) *inputoutput.ReviewApplicationInput {
				return &inputoutput.ReviewApplicationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, ApplicationID: app.ID, Status: string(constants.ApplicationStatusApproved)}
			},
			wantErr: domain.ErrorTypeForbidden,
		},
		{
			name: "application already reviewed",
			setup: func(pr *mocks.ProjectRepository, ar *mocks.ProjectApplicationRepository, _ *mocks.ProjectMemberRepository, profRepo *mocks.ProfessorRepository, p *entities.Project, prof *entities.Professor, app *entities.ProjectApplication) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return prof, nil }
				app.Status = constants.ApplicationStatusApproved
				ar.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }
			},
			input: func(p *entities.Project, prof *entities.Professor, app *entities.ProjectApplication) *inputoutput.ReviewApplicationInput {
				return &inputoutput.ReviewApplicationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, ApplicationID: app.ID, Status: string(constants.ApplicationStatusApproved)}
			},
			wantErr: domain.ErrorTypeConflict,
		},
		{
			name: "approve when slots full",
			setup: func(pr *mocks.ProjectRepository, ar *mocks.ProjectApplicationRepository, mr *mocks.ProjectMemberRepository, profRepo *mocks.ProfessorRepository, p *entities.Project, prof *entities.Professor, app *entities.ProjectApplication) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return prof, nil }
				ar.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }
				mr.CountActiveFn = func(_ context.Context, _ string) (int, error) { return p.Slots, nil }
			},
			input: func(p *entities.Project, prof *entities.Professor, app *entities.ProjectApplication) *inputoutput.ReviewApplicationInput {
				return &inputoutput.ReviewApplicationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, ApplicationID: app.ID, Status: string(constants.ApplicationStatusApproved)}
			},
			wantErr: domain.ErrorTypeConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			project, professor, app := reviewFixtures()
			projectRepo := &mocks.ProjectRepository{}
			appRepo := &mocks.ProjectApplicationRepository{}
			memberRepo := &mocks.ProjectMemberRepository{}
			profRepo := &mocks.ProfessorRepository{}
			tc.setup(projectRepo, appRepo, memberRepo, profRepo, project, professor, app)

			uc := professor_usecase.NewReviewApplicationUsecase(projectRepo, appRepo, memberRepo, profRepo)
			_, err := uc.Execute(context.Background(), tc.input(project, professor, app))
			utils.AssertAppErr(t, err, tc.wantErr)
		})
	}
}

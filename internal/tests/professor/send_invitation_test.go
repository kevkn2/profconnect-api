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

func sendInvitationFixtures() (project *entities.Project, professor *entities.Professor, student *entities.Student) {
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
	student = &entities.Student{
		ID:   "stu-1",
		User: &entities.User{ID: "user-stu-1", Name: "Stu", Email: "stu@example.com"},
	}
	return
}

func newSendInvitationDeps(t *testing.T) (
	*mocks.ProjectRepository,
	*mocks.ProjectInvitationRepository,
	*mocks.ProjectMemberRepository,
	*mocks.ProjectApplicationRepository,
	*mocks.ProfessorRepository,
	*mocks.StudentsRepository,
) {
	t.Helper()
	return &mocks.ProjectRepository{},
		&mocks.ProjectInvitationRepository{},
		&mocks.ProjectMemberRepository{},
		&mocks.ProjectApplicationRepository{},
		&mocks.ProfessorRepository{},
		&mocks.StudentsRepository{}
}

func TestSendInvitation_Success(t *testing.T) {
	project, professor, student := sendInvitationFixtures()
	projectRepo, invRepo, memberRepo, appRepo, profRepo, stuRepo := newSendInvitationDeps(t)

	projectRepo.GetByIDFn = func(_ context.Context, id string) (*entities.Project, error) {
		if id != project.ID {
			t.Fatalf("got project id %q", id)
		}
		return project, nil
	}
	profRepo.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
	stuRepo.GetStudentByIDFn = func(_ context.Context, id string) (*entities.Student, error) {
		if id != student.ID {
			t.Fatalf("got student id %q", id)
		}
		return student, nil
	}
	memberRepo.GetActiveByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectMember, error) {
		return nil, nil
	}
	appRepo.GetByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectApplication, error) {
		return nil, nil
	}
	invRepo.GetByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectInvitation, error) {
		return nil, nil
	}
	createCalled := false
	invRepo.CreateFn = func(_ context.Context, inv *entities.ProjectInvitation) (*entities.ProjectInvitation, error) {
		createCalled = true
		if inv.Status != constants.InvitationStatusPending {
			t.Fatalf("expected pending status, got %q", inv.Status)
		}
		if inv.Student.ID != student.ID {
			t.Fatalf("expected student id %q, got %q", student.ID, inv.Student.ID)
		}
		return &entities.ProjectInvitation{
			ID:      "inv-1",
			Project: inv.Project,
			Student: inv.Student,
			Status:  inv.Status,
			Message: inv.Message,
		}, nil
	}

	uc := professor_usecase.NewSendInvitationUsecase(projectRepo, invRepo, memberRepo, appRepo, profRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.SendInvitationInput{
		ProfessorUserID: professor.User.ID,
		ProjectID:       project.ID,
		StudentID:       student.ID,
		Message:         "Join us",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !createCalled {
		t.Fatalf("expected invitation Create to be called")
	}
	if out.ID != "inv-1" || out.Status != string(constants.InvitationStatusPending) {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestSendInvitation_ValidationAndPreconditions(t *testing.T) {
	cases := []struct {
		name    string
		setup   func(*mocks.ProjectRepository, *mocks.ProjectInvitationRepository, *mocks.ProjectMemberRepository, *mocks.ProjectApplicationRepository, *mocks.ProfessorRepository, *mocks.StudentsRepository, *entities.Project, *entities.Professor, *entities.Student)
		input   func(*entities.Project, *entities.Professor, *entities.Student) *inputoutput.SendInvitationInput
		wantErr domain.ErrorType
	}{
		{
			name: "missing student id",
			setup: func(_ *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, _ *mocks.ProjectApplicationRepository, _ *mocks.ProfessorRepository, _ *mocks.StudentsRepository, _ *entities.Project, _ *entities.Professor, _ *entities.Student) {
			},
			input: func(p *entities.Project, prof *entities.Professor, _ *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID}
			},
			wantErr: domain.ErrorTypeValidation,
		},
		{
			name: "project not found",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, _ *mocks.ProjectApplicationRepository, _ *mocks.ProfessorRepository, _ *mocks.StudentsRepository, _ *entities.Project, _ *entities.Professor, _ *entities.Student) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return nil, nil }
			},
			input: func(p *entities.Project, prof *entities.Professor, s *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, StudentID: s.ID}
			},
			wantErr: domain.ErrorTypeNotFound,
		},
		{
			name: "project closed",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, _ *mocks.ProjectApplicationRepository, _ *mocks.ProfessorRepository, _ *mocks.StudentsRepository, p *entities.Project, _ *entities.Professor, _ *entities.Student) {
				p.Status = constants.ProjectStatusClosed
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
			},
			input: func(p *entities.Project, prof *entities.Professor, s *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, StudentID: s.ID}
			},
			wantErr: domain.ErrorTypeConflict,
		},
		{
			name: "not project owner",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, _ *mocks.ProjectApplicationRepository, prof *mocks.ProfessorRepository, _ *mocks.StudentsRepository, p *entities.Project, _ *entities.Professor, _ *entities.Student) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				prof.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) {
					return &entities.Professor{ID: "prof-other"}, nil
				}
			},
			input: func(p *entities.Project, prof *entities.Professor, s *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, StudentID: s.ID}
			},
			wantErr: domain.ErrorTypeForbidden,
		},
		{
			name: "student not found",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, _ *mocks.ProjectApplicationRepository, prof *mocks.ProfessorRepository, stu *mocks.StudentsRepository, p *entities.Project, professor *entities.Professor, _ *entities.Student) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				prof.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
				stu.GetStudentByIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return nil, nil }
			},
			input: func(p *entities.Project, prof *entities.Professor, s *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, StudentID: s.ID}
			},
			wantErr: domain.ErrorTypeNotFound,
		},
		{
			name: "already a member",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, mem *mocks.ProjectMemberRepository, _ *mocks.ProjectApplicationRepository, prof *mocks.ProfessorRepository, stu *mocks.StudentsRepository, p *entities.Project, professor *entities.Professor, s *entities.Student) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				prof.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
				stu.GetStudentByIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return s, nil }
				mem.GetActiveByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectMember, error) {
					return &entities.ProjectMember{ID: "mem-1"}, nil
				}
			},
			input: func(p *entities.Project, prof *entities.Professor, s *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, StudentID: s.ID}
			},
			wantErr: domain.ErrorTypeConflict,
		},
		{
			name: "pending application exists",
			setup: func(pr *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, mem *mocks.ProjectMemberRepository, app *mocks.ProjectApplicationRepository, prof *mocks.ProfessorRepository, stu *mocks.StudentsRepository, p *entities.Project, professor *entities.Professor, s *entities.Student) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				prof.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
				stu.GetStudentByIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return s, nil }
				mem.GetActiveByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectMember, error) { return nil, nil }
				app.GetByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectApplication, error) {
					return &entities.ProjectApplication{ID: "app-1", Status: constants.ApplicationStatusPending}, nil
				}
			},
			input: func(p *entities.Project, prof *entities.Professor, s *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, StudentID: s.ID}
			},
			wantErr: domain.ErrorTypeConflict,
		},
		{
			name: "pending invitation already exists",
			setup: func(pr *mocks.ProjectRepository, inv *mocks.ProjectInvitationRepository, mem *mocks.ProjectMemberRepository, app *mocks.ProjectApplicationRepository, prof *mocks.ProfessorRepository, stu *mocks.StudentsRepository, p *entities.Project, professor *entities.Professor, s *entities.Student) {
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				prof.GetByUserIDFn = func(_ context.Context, _ string) (*entities.Professor, error) { return professor, nil }
				stu.GetStudentByIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return s, nil }
				mem.GetActiveByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectMember, error) { return nil, nil }
				app.GetByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectApplication, error) { return nil, nil }
				inv.GetByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectInvitation, error) {
					return &entities.ProjectInvitation{ID: "inv-1", Status: constants.InvitationStatusPending}, nil
				}
			},
			input: func(p *entities.Project, prof *entities.Professor, s *entities.Student) *inputoutput.SendInvitationInput {
				return &inputoutput.SendInvitationInput{ProfessorUserID: prof.User.ID, ProjectID: p.ID, StudentID: s.ID}
			},
			wantErr: domain.ErrorTypeConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			project, professor, student := sendInvitationFixtures()
			projectRepo, invRepo, memberRepo, appRepo, profRepo, stuRepo := newSendInvitationDeps(t)
			tc.setup(projectRepo, invRepo, memberRepo, appRepo, profRepo, stuRepo, project, professor, student)

			uc := professor_usecase.NewSendInvitationUsecase(projectRepo, invRepo, memberRepo, appRepo, profRepo, stuRepo)
			_, err := uc.Execute(context.Background(), tc.input(project, professor, student))
			utils.AssertAppErr(t, err, tc.wantErr)
		})
	}
}

func TestSendInvitation_PropagatesRepoErrors(t *testing.T) {
	project, professor, _ := sendInvitationFixtures()
	projectRepo, invRepo, memberRepo, appRepo, profRepo, stuRepo := newSendInvitationDeps(t)

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) {
		return nil, errors.New("db down")
	}

	uc := professor_usecase.NewSendInvitationUsecase(projectRepo, invRepo, memberRepo, appRepo, profRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.SendInvitationInput{
		ProfessorUserID: professor.User.ID,
		ProjectID:       project.ID,
		StudentID:       "stu-1",
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

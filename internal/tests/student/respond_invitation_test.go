package student

import (
	"context"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	student_usecase "profconnect-api/internal/usecase/student"
)

func respondFixtures() (project *entities.Project, student *entities.Student, invitation *entities.ProjectInvitation) {
	student = &entities.Student{
		ID:   "stu-1",
		User: &entities.User{ID: "user-stu-1", Name: "Stu", Email: "stu@example.com"},
	}
	project = &entities.Project{
		ID:     "proj-1",
		Title:  "AI Research",
		Slots:  2,
		Status: constants.ProjectStatusOpen,
	}
	invitation = &entities.ProjectInvitation{
		ID:      "inv-1",
		Project: &entities.Project{ID: project.ID},
		Student: &entities.Student{ID: student.ID},
		Status:  constants.InvitationStatusPending,
	}
	return
}

func TestRespondInvitation_AcceptCreatesMember(t *testing.T) {
	project, student, invitation := respondFixtures()
	projectRepo := &mocks.ProjectRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	stuRepo := &mocks.StudentsRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	invRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return invitation, nil }
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	memberRepo.CountActiveFn = func(_ context.Context, _ string) (int, error) { return 1, nil }
	invRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.InvitationStatus) (*entities.ProjectInvitation, error) {
		if s != constants.InvitationStatusAccepted {
			t.Fatalf("expected accepted, got %q", s)
		}
		return &entities.ProjectInvitation{
			ID:      invitation.ID,
			Project: invitation.Project,
			Student: invitation.Student,
			Status:  s,
		}, nil
	}
	createCalled := false
	memberRepo.CreateFn = func(_ context.Context, m *entities.ProjectMember) (*entities.ProjectMember, error) {
		createCalled = true
		if m.Source != constants.MemberSourceInvitation {
			t.Fatalf("expected source 'invitation', got %q", m.Source)
		}
		if m.SourceRefID != invitation.ID {
			t.Fatalf("expected source ref %q, got %q", invitation.ID, m.SourceRefID)
		}
		return m, nil
	}

	uc := student_usecase.NewRespondInvitationUsecase(projectRepo, invRepo, memberRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.RespondInvitationInput{
		StudentUserID: student.User.ID,
		InvitationID:  invitation.ID,
		Status:        string(constants.InvitationStatusAccepted),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !createCalled {
		t.Fatalf("expected member Create to be called")
	}
	if out.Status != string(constants.InvitationStatusAccepted) {
		t.Fatalf("expected accepted, got %q", out.Status)
	}
}

func TestRespondInvitation_AcceptClosesProjectWhenSlotsFill(t *testing.T) {
	project, student, invitation := respondFixtures()
	project.Slots = 1
	projectRepo := &mocks.ProjectRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	stuRepo := &mocks.StudentsRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	invRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return invitation, nil }
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	calls := 0
	memberRepo.CountActiveFn = func(_ context.Context, _ string) (int, error) {
		calls++
		if calls == 1 {
			return 0, nil
		}
		return 1, nil
	}
	invRepo.UpdateStatusFn = func(_ context.Context, _ string, _ constants.InvitationStatus) (*entities.ProjectInvitation, error) {
		return &entities.ProjectInvitation{ID: invitation.ID, Project: invitation.Project, Student: invitation.Student, Status: constants.InvitationStatusAccepted}, nil
	}
	memberRepo.CreateFn = func(_ context.Context, m *entities.ProjectMember) (*entities.ProjectMember, error) { return m, nil }
	closed := false
	projectRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.ProjectStatus) error {
		if s == constants.ProjectStatusClosed {
			closed = true
		}
		return nil
	}

	uc := student_usecase.NewRespondInvitationUsecase(projectRepo, invRepo, memberRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.RespondInvitationInput{
		StudentUserID: student.User.ID,
		InvitationID:  invitation.ID,
		Status:        string(constants.InvitationStatusAccepted),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !closed {
		t.Fatalf("expected project to be closed when slots fill")
	}
}

func TestRespondInvitation_DeclineDoesNotCreateMember(t *testing.T) {
	project, student, invitation := respondFixtures()
	projectRepo := &mocks.ProjectRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	stuRepo := &mocks.StudentsRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	invRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return invitation, nil }
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	invRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.InvitationStatus) (*entities.ProjectInvitation, error) {
		if s != constants.InvitationStatusDeclined {
			t.Fatalf("expected declined, got %q", s)
		}
		return &entities.ProjectInvitation{ID: invitation.ID, Project: invitation.Project, Student: invitation.Student, Status: s}, nil
	}

	uc := student_usecase.NewRespondInvitationUsecase(projectRepo, invRepo, memberRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.RespondInvitationInput{
		StudentUserID: student.User.ID,
		InvitationID:  invitation.ID,
		Status:        string(constants.InvitationStatusDeclined),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != string(constants.InvitationStatusDeclined) {
		t.Fatalf("expected declined, got %q", out.Status)
	}
}

func TestRespondInvitation_ValidationAndPreconditions(t *testing.T) {
	cases := []struct {
		name    string
		setup   func(*mocks.ProjectRepository, *mocks.ProjectInvitationRepository, *mocks.ProjectMemberRepository, *mocks.StudentsRepository, *entities.Project, *entities.Student, *entities.ProjectInvitation)
		input   func(*entities.Student, *entities.ProjectInvitation) *inputoutput.RespondInvitationInput
		wantErr domain.ErrorType
	}{
		{
			name: "invalid status",
			setup: func(_ *mocks.ProjectRepository, _ *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, _ *mocks.StudentsRepository, _ *entities.Project, _ *entities.Student, _ *entities.ProjectInvitation) {
			},
			input: func(s *entities.Student, inv *entities.ProjectInvitation) *inputoutput.RespondInvitationInput {
				return &inputoutput.RespondInvitationInput{StudentUserID: s.User.ID, InvitationID: inv.ID, Status: "expired"}
			},
			wantErr: domain.ErrorTypeValidation,
		},
		{
			name: "not addressee",
			setup: func(_ *mocks.ProjectRepository, ir *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, sr *mocks.StudentsRepository, _ *entities.Project, s *entities.Student, inv *entities.ProjectInvitation) {
				sr.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return s, nil }
				ir.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) {
					return &entities.ProjectInvitation{ID: inv.ID, Status: constants.InvitationStatusPending, Student: &entities.Student{ID: "other"}, Project: inv.Project}, nil
				}
			},
			input: func(s *entities.Student, inv *entities.ProjectInvitation) *inputoutput.RespondInvitationInput {
				return &inputoutput.RespondInvitationInput{StudentUserID: s.User.ID, InvitationID: inv.ID, Status: string(constants.InvitationStatusAccepted)}
			},
			wantErr: domain.ErrorTypeForbidden,
		},
		{
			name: "already responded",
			setup: func(_ *mocks.ProjectRepository, ir *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, sr *mocks.StudentsRepository, _ *entities.Project, s *entities.Student, inv *entities.ProjectInvitation) {
				sr.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return s, nil }
				ir.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) {
					inv.Status = constants.InvitationStatusDeclined
					return inv, nil
				}
			},
			input: func(s *entities.Student, inv *entities.ProjectInvitation) *inputoutput.RespondInvitationInput {
				return &inputoutput.RespondInvitationInput{StudentUserID: s.User.ID, InvitationID: inv.ID, Status: string(constants.InvitationStatusAccepted)}
			},
			wantErr: domain.ErrorTypeConflict,
		},
		{
			name: "accept when project closed",
			setup: func(pr *mocks.ProjectRepository, ir *mocks.ProjectInvitationRepository, _ *mocks.ProjectMemberRepository, sr *mocks.StudentsRepository, p *entities.Project, s *entities.Student, inv *entities.ProjectInvitation) {
				sr.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return s, nil }
				ir.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return inv, nil }
				p.Status = constants.ProjectStatusClosed
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
			},
			input: func(s *entities.Student, inv *entities.ProjectInvitation) *inputoutput.RespondInvitationInput {
				return &inputoutput.RespondInvitationInput{StudentUserID: s.User.ID, InvitationID: inv.ID, Status: string(constants.InvitationStatusAccepted)}
			},
			wantErr: domain.ErrorTypeConflict,
		},
		{
			name: "accept when slots filled",
			setup: func(pr *mocks.ProjectRepository, ir *mocks.ProjectInvitationRepository, mr *mocks.ProjectMemberRepository, sr *mocks.StudentsRepository, p *entities.Project, s *entities.Student, inv *entities.ProjectInvitation) {
				sr.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return s, nil }
				ir.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectInvitation, error) { return inv, nil }
				pr.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return p, nil }
				mr.CountActiveFn = func(_ context.Context, _ string) (int, error) { return p.Slots, nil }
			},
			input: func(s *entities.Student, inv *entities.ProjectInvitation) *inputoutput.RespondInvitationInput {
				return &inputoutput.RespondInvitationInput{StudentUserID: s.User.ID, InvitationID: inv.ID, Status: string(constants.InvitationStatusAccepted)}
			},
			wantErr: domain.ErrorTypeConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			project, student, invitation := respondFixtures()
			projectRepo := &mocks.ProjectRepository{}
			invRepo := &mocks.ProjectInvitationRepository{}
			memberRepo := &mocks.ProjectMemberRepository{}
			stuRepo := &mocks.StudentsRepository{}
			tc.setup(projectRepo, invRepo, memberRepo, stuRepo, project, student, invitation)

			uc := student_usecase.NewRespondInvitationUsecase(projectRepo, invRepo, memberRepo, stuRepo)
			_, err := uc.Execute(context.Background(), tc.input(student, invitation))
			utils.AssertAppErr(t, err, tc.wantErr)
		})
	}
}

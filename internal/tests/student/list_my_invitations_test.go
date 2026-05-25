package student

import (
	"context"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	student_usecase "profconnect-api/internal/usecase/student"
)

func TestListMyInvitations_Success(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	invRepo := &mocks.ProjectInvitationRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return &entities.Student{ID: "stu-1"}, nil
	}
	invRepo.ListByStudentFn = func(_ context.Context, _ string) ([]*port.StudentInvitationView, error) {
		return []*port.StudentInvitationView{
			{ID: "i-1", Status: constants.InvitationStatusPending, ProjectTitle: "P1"},
			{ID: "i-2", Status: constants.InvitationStatusDeclined, ProjectTitle: "P2"},
		}, nil
	}

	uc := student_usecase.NewListMyInvitationsUsecase(invRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListMyInvitationsInput{StudentUserID: "u"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Invitations) != 2 {
		t.Fatalf("expected 2 invitations, got %d", len(out.Invitations))
	}
	if out.Invitations[0].Project.Title != "P1" {
		t.Fatalf("expected project title propagated from view")
	}
}

func TestListMyInvitations_StudentNotFound(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return nil, nil }

	uc := student_usecase.NewListMyInvitationsUsecase(&mocks.ProjectInvitationRepository{}, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ListMyInvitationsInput{})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

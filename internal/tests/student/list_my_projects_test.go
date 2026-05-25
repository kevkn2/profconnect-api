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

func TestListMyProjects_Success(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return &entities.Student{ID: "stu-1"}, nil
	}
	memberRepo.ListActiveByStudentFn = func(_ context.Context, _ string) ([]*port.StudentMembershipView, error) {
		return []*port.StudentMembershipView{
			{ID: "m-1", Status: constants.MemberStatusActive, Source: constants.MemberSourceApplication, ProjectTitle: "P1"},
			{ID: "m-2", Status: constants.MemberStatusActive, Source: constants.MemberSourceInvitation, ProjectTitle: "P2"},
		}, nil
	}

	uc := student_usecase.NewListMyProjectsUsecase(memberRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListMyProjectsInput{StudentUserID: "u"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Memberships) != 2 {
		t.Fatalf("expected 2 memberships, got %d", len(out.Memberships))
	}
	if out.Memberships[0].Source != string(constants.MemberSourceApplication) {
		t.Fatalf("unexpected first source: %q", out.Memberships[0].Source)
	}
}

func TestListMyProjects_StudentNotFound(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return nil, nil }

	uc := student_usecase.NewListMyProjectsUsecase(&mocks.ProjectMemberRepository{}, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ListMyProjectsInput{})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

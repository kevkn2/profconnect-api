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

func TestMyApplications_Success(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return &entities.Student{ID: "stu-1"}, nil
	}
	appRepo.ListByStudentFn = func(_ context.Context, _ string) ([]*port.StudentApplicationView, error) {
		return []*port.StudentApplicationView{
			{ID: "a-1", Status: constants.ApplicationStatusPending, ProjectTitle: "P1"},
			{ID: "a-2", Status: constants.ApplicationStatusApproved, ProjectTitle: "P2"},
		}, nil
	}

	uc := student_usecase.NewListMyApplicationsUsecase(appRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListMyApplicationsInput{StudentUserID: "u"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Applications) != 2 {
		t.Fatalf("expected 2 applications, got %d", len(out.Applications))
	}
	if out.Applications[1].Project.Title != "P2" {
		t.Fatalf("expected project title propagated from view")
	}
}

func TestMyApplications_StudentNotFound(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return nil, nil }

	uc := student_usecase.NewListMyApplicationsUsecase(&mocks.ProjectApplicationRepository{}, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ListMyApplicationsInput{})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

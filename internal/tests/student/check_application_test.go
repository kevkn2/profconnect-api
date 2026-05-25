package student

import (
	"context"
	"errors"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	test_utils "profconnect-api/internal/tests/utils"
	student_usecase "profconnect-api/internal/usecase/student"
)

func TestCheckApplication_Exists(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return &entities.Student{ID: "stu-1"}, nil
	}
	appRepo.CheckApplicationStatusFn = func(_ context.Context, _, _ string) (bool, error) { return true, nil }

	uc := student_usecase.NewListApplicationsPerIDUsecase(appRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.CheckApplicationStatusInput{
		StudentUserID: "u", ProjectID: "p",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Exists {
		t.Fatalf("expected Exists=true")
	}
}

func TestCheckApplication_DoesNotExist(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return &entities.Student{ID: "stu-1"}, nil
	}
	appRepo.CheckApplicationStatusFn = func(_ context.Context, _, _ string) (bool, error) { return false, nil }

	uc := student_usecase.NewListApplicationsPerIDUsecase(appRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.CheckApplicationStatusInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Exists {
		t.Fatalf("expected Exists=false")
	}
}

func TestCheckApplication_StudentNotFound(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return nil, nil }

	uc := student_usecase.NewListApplicationsPerIDUsecase(&mocks.ProjectApplicationRepository{}, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CheckApplicationStatusInput{})
	test_utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

func TestCheckApplication_RepoError(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return &entities.Student{ID: "stu-1"}, nil
	}
	appRepo.CheckApplicationStatusFn = func(_ context.Context, _, _ string) (bool, error) {
		return false, errors.New("db down")
	}

	uc := student_usecase.NewListApplicationsPerIDUsecase(appRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.CheckApplicationStatusInput{})
	test_utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

package student

import (
	"context"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	test_utils "profconnect-api/internal/tests/utils"
	student_usecase "profconnect-api/internal/usecase/student"
)

func applyFixtures() (project *entities.Project, student *entities.Student) {
	student = &entities.Student{
		ID:   "stu-1",
		User: &entities.User{ID: "user-stu-1", Name: "Stu", Email: "stu@x"},
	}
	project = &entities.Project{ID: "proj-1", Title: "AI", Description: "d", Slots: 2, Status: constants.ProjectStatusOpen}
	return
}

func TestApplyProject_Success(t *testing.T) {
	project, student := applyFixtures()
	projectRepo := &mocks.ProjectRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	stuRepo := &mocks.StudentsRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	appRepo.GetByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectApplication, error) { return nil, nil }
	created := false
	appRepo.CreateFn = func(_ context.Context, a *entities.ProjectApplication) (*entities.ProjectApplication, error) {
		created = true
		if a.Status != constants.ApplicationStatusPending {
			t.Fatalf("expected pending, got %q", a.Status)
		}
		return &entities.ProjectApplication{
			ID: "app-1", Project: a.Project, Student: a.Student, Status: a.Status, Message: a.Message,
		}, nil
	}

	uc := student_usecase.NewApplyProjectUsecase(projectRepo, appRepo, stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ApplyProjectInput{
		StudentUserID: student.User.ID, ProjectID: project.ID, Message: "hi",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatalf("expected Create to be called")
	}
	if out.ID != "app-1" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestApplyProject_ProjectClosed(t *testing.T) {
	project, student := applyFixtures()
	project.Status = constants.ProjectStatusClosed
	projectRepo := &mocks.ProjectRepository{}
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }

	uc := student_usecase.NewApplyProjectUsecase(projectRepo, &mocks.ProjectApplicationRepository{}, &mocks.StudentsRepository{})
	_, err := uc.Execute(context.Background(), &inputoutput.ApplyProjectInput{
		StudentUserID: student.User.ID, ProjectID: project.ID,
	})
	test_utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestApplyProject_AlreadyApplied(t *testing.T) {
	project, student := applyFixtures()
	projectRepo := &mocks.ProjectRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	stuRepo := &mocks.StudentsRepository{}

	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	appRepo.GetByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectApplication, error) {
		return &entities.ProjectApplication{ID: "existing"}, nil
	}

	uc := student_usecase.NewApplyProjectUsecase(projectRepo, appRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ApplyProjectInput{
		StudentUserID: student.User.ID, ProjectID: project.ID,
	})
	test_utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestApplyProject_ProjectNotFound(t *testing.T) {
	projectRepo := &mocks.ProjectRepository{}
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return nil, nil }

	uc := student_usecase.NewApplyProjectUsecase(projectRepo, &mocks.ProjectApplicationRepository{}, &mocks.StudentsRepository{})
	_, err := uc.Execute(context.Background(), &inputoutput.ApplyProjectInput{ProjectID: "x"})
	test_utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

func TestApplyProject_StudentNotFound(t *testing.T) {
	project, _ := applyFixtures()
	projectRepo := &mocks.ProjectRepository{}
	stuRepo := &mocks.StudentsRepository{}
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return nil, nil }

	uc := student_usecase.NewApplyProjectUsecase(projectRepo, &mocks.ProjectApplicationRepository{}, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ApplyProjectInput{ProjectID: project.ID})
	test_utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

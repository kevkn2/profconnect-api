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

func withdrawFixtures() (student *entities.Student, application *entities.ProjectApplication) {
	student = &entities.Student{ID: "stu-1", User: &entities.User{ID: "user-stu-1"}}
	application = &entities.ProjectApplication{
		ID:      "app-1",
		Project: &entities.Project{ID: "proj-1"},
		Student: &entities.Student{ID: student.ID},
		Status:  constants.ApplicationStatusPending,
	}
	return
}

func TestWithdrawApplication_Success(t *testing.T) {
	student, app := withdrawFixtures()
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }
	deleted := false
	appRepo.DeleteFn = func(_ context.Context, id string) error {
		deleted = true
		if id != app.ID {
			t.Fatalf("expected delete on %q, got %q", app.ID, id)
		}
		return nil
	}

	uc := student_usecase.NewWithdrawApplicationUsecase(appRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.WithdrawApplicationInput{
		StudentUserID: student.User.ID, ProjectID: app.Project.ID, ApplicationID: app.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Fatalf("expected Delete to be called")
	}
}

func TestWithdrawApplication_ApprovedCannotWithdraw(t *testing.T) {
	student, app := withdrawFixtures()
	app.Status = constants.ApplicationStatusApproved
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }

	uc := student_usecase.NewWithdrawApplicationUsecase(appRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.WithdrawApplicationInput{
		StudentUserID: student.User.ID, ProjectID: app.Project.ID, ApplicationID: app.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeConflict)
}

func TestWithdrawApplication_NotOwner(t *testing.T) {
	student, app := withdrawFixtures()
	app.Student = &entities.Student{ID: "other"}
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }

	uc := student_usecase.NewWithdrawApplicationUsecase(appRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.WithdrawApplicationInput{
		StudentUserID: student.User.ID, ProjectID: app.Project.ID, ApplicationID: app.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeForbidden)
}

func TestWithdrawApplication_WrongProject(t *testing.T) {
	student, app := withdrawFixtures()
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return app, nil }

	uc := student_usecase.NewWithdrawApplicationUsecase(appRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.WithdrawApplicationInput{
		StudentUserID: student.User.ID, ProjectID: "wrong-proj", ApplicationID: app.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

func TestWithdrawApplication_ApplicationNotFound(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	appRepo := &mocks.ProjectApplicationRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return &entities.Student{ID: "s"}, nil
	}
	appRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.ProjectApplication, error) { return nil, nil }

	uc := student_usecase.NewWithdrawApplicationUsecase(appRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.WithdrawApplicationInput{ApplicationID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

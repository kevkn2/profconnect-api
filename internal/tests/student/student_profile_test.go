package student

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
	student_usecase "profconnect-api/internal/usecase/student"
)

func TestStudentProfile_Success(t *testing.T) {
	student := &entities.Student{
		ID:                "stu-1",
		User:              &entities.User{ID: "user-1", Name: "Stu", Email: "s@x", Role: constants.Student},
		University:        "MIT",
		Department:        "CS",
		ResearchInterests: "AI",
	}
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }

	uc := student_usecase.NewProfileUsecase(stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ProfileInput{UserID: student.User.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != student.ID || out.University != "MIT" || out.ResearchInterests != "AI" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestStudentProfile_NotFound(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return nil, nil }

	uc := student_usecase.NewProfileUsecase(stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ProfileInput{UserID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

func TestStudentProfile_DefaultsRoleWhenEmpty(t *testing.T) {
	student := &entities.Student{
		ID:   "stu-1",
		User: &entities.User{ID: "user-1", Name: "Stu", Email: "s@x"}, // role left empty
	}
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }

	uc := student_usecase.NewProfileUsecase(stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ProfileInput{UserID: student.User.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Role != string(constants.Student) {
		t.Fatalf("expected role to default to student, got %q", out.Role)
	}
}

func TestStudentProfile_RepoError(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) {
		return nil, errors.New("db down")
	}

	uc := student_usecase.NewProfileUsecase(stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ProfileInput{UserID: "x"})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

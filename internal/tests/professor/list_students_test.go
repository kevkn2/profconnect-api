package professor

import (
	"context"
	"errors"
	"testing"

	"profconnect-api/internal/domain"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/tests/mocks"
	"profconnect-api/internal/tests/utils"
	professor_usecase "profconnect-api/internal/usecase/professor"
)

func TestListStudents_Success(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.ListStudentsFn = func(_ context.Context) ([]*entities.Student, error) {
		return []*entities.Student{
			{ID: "s-1", User: &entities.User{ID: "u-1", Name: "A", Email: "a@x"}, University: "MIT", Department: "CS"},
			{ID: "s-2", User: &entities.User{ID: "u-2", Name: "B", Email: "b@x"}, University: "Stanford", Department: "EE"},
		}, nil
	}

	uc := professor_usecase.NewListStudentsUsecase(stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListStudentsInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Students) != 2 {
		t.Fatalf("expected 2 students, got %d", len(out.Students))
	}
	if out.Students[0].Name != "A" || out.Students[0].UserID != "u-1" {
		t.Fatalf("unexpected first student: %+v", out.Students[0])
	}
}

func TestListStudents_RepoError(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.ListStudentsFn = func(_ context.Context) ([]*entities.Student, error) {
		return nil, errors.New("db down")
	}

	uc := professor_usecase.NewListStudentsUsecase(stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.ListStudentsInput{})
	utils.AssertAppErr(t, err, domain.ErrorTypeInternal)
}

func TestListStudents_Empty(t *testing.T) {
	stuRepo := &mocks.StudentsRepository{}
	stuRepo.ListStudentsFn = func(_ context.Context) ([]*entities.Student, error) { return nil, nil }

	uc := professor_usecase.NewListStudentsUsecase(stuRepo)
	out, err := uc.Execute(context.Background(), &inputoutput.ListStudentsInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Students) != 0 {
		t.Fatalf("expected empty list, got %d", len(out.Students))
	}
}

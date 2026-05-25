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

func leaveFixtures() (project *entities.Project, student *entities.Student, member *entities.ProjectMember) {
	student = &entities.Student{ID: "stu-1", User: &entities.User{ID: "user-stu-1"}}
	project = &entities.Project{ID: "proj-1", Slots: 2, Status: constants.ProjectStatusOpen}
	member = &entities.ProjectMember{ID: "mem-1", Status: constants.MemberStatusActive}
	return
}

func TestLeaveProject_Success(t *testing.T) {
	project, student, member := leaveFixtures()
	projectRepo := &mocks.ProjectRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	stuRepo := &mocks.StudentsRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	memberRepo.GetActiveByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectMember, error) {
		return member, nil
	}
	updated := false
	memberRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.MemberStatus) (*entities.ProjectMember, error) {
		updated = true
		if s != constants.MemberStatusLeft {
			t.Fatalf("expected status 'left', got %q", s)
		}
		return member, nil
	}

	uc := student_usecase.NewLeaveProjectUsecase(projectRepo, memberRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.LeaveProjectInput{
		StudentUserID: student.User.ID,
		ProjectID:     project.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated {
		t.Fatalf("expected UpdateStatus to be called")
	}
}

func TestLeaveProject_ReopensClosedProject(t *testing.T) {
	project, student, member := leaveFixtures()
	project.Status = constants.ProjectStatusClosed
	projectRepo := &mocks.ProjectRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	stuRepo := &mocks.StudentsRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	memberRepo.GetActiveByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectMember, error) {
		return member, nil
	}
	memberRepo.UpdateStatusFn = func(_ context.Context, _ string, _ constants.MemberStatus) (*entities.ProjectMember, error) {
		return member, nil
	}
	reopened := false
	projectRepo.UpdateStatusFn = func(_ context.Context, _ string, s constants.ProjectStatus) error {
		if s == constants.ProjectStatusOpen {
			reopened = true
		}
		return nil
	}

	uc := student_usecase.NewLeaveProjectUsecase(projectRepo, memberRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.LeaveProjectInput{
		StudentUserID: student.User.ID,
		ProjectID:     project.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reopened {
		t.Fatalf("expected closed project to be reopened when slot frees")
	}
}

func TestLeaveProject_NotAMember(t *testing.T) {
	project, student, _ := leaveFixtures()
	projectRepo := &mocks.ProjectRepository{}
	memberRepo := &mocks.ProjectMemberRepository{}
	stuRepo := &mocks.StudentsRepository{}

	stuRepo.GetStudentByUserIDFn = func(_ context.Context, _ string) (*entities.Student, error) { return student, nil }
	projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) { return project, nil }
	memberRepo.GetActiveByProjectAndStudentFn = func(_ context.Context, _, _ string) (*entities.ProjectMember, error) {
		return nil, nil
	}

	uc := student_usecase.NewLeaveProjectUsecase(projectRepo, memberRepo, stuRepo)
	_, err := uc.Execute(context.Background(), &inputoutput.LeaveProjectInput{
		StudentUserID: student.User.ID,
		ProjectID:     project.ID,
	})
	utils.AssertAppErr(t, err, domain.ErrorTypeNotFound)
}

// Package mocks provides test doubles for the domain port interfaces.
// Each mock exposes function fields so individual tests can override only the
// methods they exercise; un-stubbed methods panic to make accidental usage
// visible.
package mocks

import (
	"context"
	"fmt"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/domain/port"
)

func unimplemented(method string) error {
	return fmt.Errorf("mock: %s not implemented", method)
}

// ProjectRepository is a mock of port.ProjectRepository.
type ProjectRepository struct {
	CreateFn          func(ctx context.Context, p *entities.Project) (*entities.Project, error)
	GetByIDFn         func(ctx context.Context, id string) (*entities.Project, error)
	ListFn            func(ctx context.Context) ([]*entities.Project, error)
	ListByProfessorFn func(ctx context.Context, professorID string) ([]*entities.Project, error)
	UpdateStatusFn    func(ctx context.Context, id string, status constants.ProjectStatus) error
}

func (m *ProjectRepository) Create(ctx context.Context, p *entities.Project) (*entities.Project, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, p)
	}
	return nil, unimplemented("ProjectRepository.Create")
}
func (m *ProjectRepository) GetByID(ctx context.Context, id string) (*entities.Project, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, unimplemented("ProjectRepository.GetByID")
}
func (m *ProjectRepository) List(ctx context.Context) ([]*entities.Project, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}
	return nil, unimplemented("ProjectRepository.List")
}
func (m *ProjectRepository) ListByProfessor(ctx context.Context, professorID string) ([]*entities.Project, error) {
	if m.ListByProfessorFn != nil {
		return m.ListByProfessorFn(ctx, professorID)
	}
	return nil, unimplemented("ProjectRepository.ListByProfessor")
}
func (m *ProjectRepository) UpdateStatus(ctx context.Context, id string, status constants.ProjectStatus) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, status)
	}
	return unimplemented("ProjectRepository.UpdateStatus")
}

// ProjectApplicationRepository is a mock of port.ProjectApplicationRepository.
type ProjectApplicationRepository struct {
	CreateFn                  func(ctx context.Context, a *entities.ProjectApplication) (*entities.ProjectApplication, error)
	GetByIDFn                 func(ctx context.Context, id string) (*entities.ProjectApplication, error)
	GetByProjectAndStudentFn  func(ctx context.Context, projectID, studentID string) (*entities.ProjectApplication, error)
	ListByProjectFn           func(ctx context.Context, projectID string) ([]*entities.ProjectApplication, error)
	ListByStudentFn           func(ctx context.Context, studentID string) ([]*port.StudentApplicationView, error)
	CheckApplicationStatusFn  func(ctx context.Context, studentID, projectID string) (bool, error)
	UpdateStatusFn            func(ctx context.Context, id string, status constants.ApplicationStatus) (*entities.ProjectApplication, error)
	DeleteFn                  func(ctx context.Context, id string) error
}

func (m *ProjectApplicationRepository) Create(ctx context.Context, a *entities.ProjectApplication) (*entities.ProjectApplication, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, a)
	}
	return nil, unimplemented("ProjectApplicationRepository.Create")
}
func (m *ProjectApplicationRepository) GetByID(ctx context.Context, id string) (*entities.ProjectApplication, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, unimplemented("ProjectApplicationRepository.GetByID")
}
func (m *ProjectApplicationRepository) GetByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectApplication, error) {
	if m.GetByProjectAndStudentFn != nil {
		return m.GetByProjectAndStudentFn(ctx, projectID, studentID)
	}
	return nil, unimplemented("ProjectApplicationRepository.GetByProjectAndStudent")
}
func (m *ProjectApplicationRepository) ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectApplication, error) {
	if m.ListByProjectFn != nil {
		return m.ListByProjectFn(ctx, projectID)
	}
	return nil, unimplemented("ProjectApplicationRepository.ListByProject")
}
func (m *ProjectApplicationRepository) ListByStudent(ctx context.Context, studentID string) ([]*port.StudentApplicationView, error) {
	if m.ListByStudentFn != nil {
		return m.ListByStudentFn(ctx, studentID)
	}
	return nil, unimplemented("ProjectApplicationRepository.ListByStudent")
}
func (m *ProjectApplicationRepository) CheckApplicationStatus(ctx context.Context, studentID, projectID string) (bool, error) {
	if m.CheckApplicationStatusFn != nil {
		return m.CheckApplicationStatusFn(ctx, studentID, projectID)
	}
	return false, unimplemented("ProjectApplicationRepository.CheckApplicationStatus")
}
func (m *ProjectApplicationRepository) UpdateStatus(ctx context.Context, id string, status constants.ApplicationStatus) (*entities.ProjectApplication, error) {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, status)
	}
	return nil, unimplemented("ProjectApplicationRepository.UpdateStatus")
}
func (m *ProjectApplicationRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return unimplemented("ProjectApplicationRepository.Delete")
}

// ProjectMemberRepository is a mock of port.ProjectMemberRepository.
type ProjectMemberRepository struct {
	CreateFn                        func(ctx context.Context, m *entities.ProjectMember) (*entities.ProjectMember, error)
	GetByIDFn                       func(ctx context.Context, id string) (*entities.ProjectMember, error)
	GetActiveByProjectAndStudentFn  func(ctx context.Context, projectID, studentID string) (*entities.ProjectMember, error)
	ListByProjectFn                 func(ctx context.Context, projectID string) ([]*entities.ProjectMember, error)
	ListActiveByStudentFn           func(ctx context.Context, studentID string) ([]*port.StudentMembershipView, error)
	CountActiveFn                   func(ctx context.Context, projectID string) (int, error)
	UpdateStatusFn                  func(ctx context.Context, id string, status constants.MemberStatus) (*entities.ProjectMember, error)
}

func (m *ProjectMemberRepository) Create(ctx context.Context, member *entities.ProjectMember) (*entities.ProjectMember, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, member)
	}
	return nil, unimplemented("ProjectMemberRepository.Create")
}
func (m *ProjectMemberRepository) GetByID(ctx context.Context, id string) (*entities.ProjectMember, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, unimplemented("ProjectMemberRepository.GetByID")
}
func (m *ProjectMemberRepository) GetActiveByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectMember, error) {
	if m.GetActiveByProjectAndStudentFn != nil {
		return m.GetActiveByProjectAndStudentFn(ctx, projectID, studentID)
	}
	return nil, unimplemented("ProjectMemberRepository.GetActiveByProjectAndStudent")
}
func (m *ProjectMemberRepository) ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectMember, error) {
	if m.ListByProjectFn != nil {
		return m.ListByProjectFn(ctx, projectID)
	}
	return nil, unimplemented("ProjectMemberRepository.ListByProject")
}
func (m *ProjectMemberRepository) ListActiveByStudent(ctx context.Context, studentID string) ([]*port.StudentMembershipView, error) {
	if m.ListActiveByStudentFn != nil {
		return m.ListActiveByStudentFn(ctx, studentID)
	}
	return nil, unimplemented("ProjectMemberRepository.ListActiveByStudent")
}
func (m *ProjectMemberRepository) CountActive(ctx context.Context, projectID string) (int, error) {
	if m.CountActiveFn != nil {
		return m.CountActiveFn(ctx, projectID)
	}
	return 0, unimplemented("ProjectMemberRepository.CountActive")
}
func (m *ProjectMemberRepository) UpdateStatus(ctx context.Context, id string, status constants.MemberStatus) (*entities.ProjectMember, error) {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, status)
	}
	return nil, unimplemented("ProjectMemberRepository.UpdateStatus")
}

// ProjectInvitationRepository is a mock of port.ProjectInvitationRepository.
type ProjectInvitationRepository struct {
	CreateFn                 func(ctx context.Context, inv *entities.ProjectInvitation) (*entities.ProjectInvitation, error)
	GetByIDFn                func(ctx context.Context, id string) (*entities.ProjectInvitation, error)
	GetByProjectAndStudentFn func(ctx context.Context, projectID, studentID string) (*entities.ProjectInvitation, error)
	ListByProjectFn          func(ctx context.Context, projectID string) ([]*entities.ProjectInvitation, error)
	ListByStudentFn          func(ctx context.Context, studentID string) ([]*port.StudentInvitationView, error)
	UpdateStatusFn           func(ctx context.Context, id string, status constants.InvitationStatus) (*entities.ProjectInvitation, error)
	DeleteFn                 func(ctx context.Context, id string) error
}

func (m *ProjectInvitationRepository) Create(ctx context.Context, inv *entities.ProjectInvitation) (*entities.ProjectInvitation, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, inv)
	}
	return nil, unimplemented("ProjectInvitationRepository.Create")
}
func (m *ProjectInvitationRepository) GetByID(ctx context.Context, id string) (*entities.ProjectInvitation, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, unimplemented("ProjectInvitationRepository.GetByID")
}
func (m *ProjectInvitationRepository) GetByProjectAndStudent(ctx context.Context, projectID, studentID string) (*entities.ProjectInvitation, error) {
	if m.GetByProjectAndStudentFn != nil {
		return m.GetByProjectAndStudentFn(ctx, projectID, studentID)
	}
	return nil, unimplemented("ProjectInvitationRepository.GetByProjectAndStudent")
}
func (m *ProjectInvitationRepository) ListByProject(ctx context.Context, projectID string) ([]*entities.ProjectInvitation, error) {
	if m.ListByProjectFn != nil {
		return m.ListByProjectFn(ctx, projectID)
	}
	return nil, unimplemented("ProjectInvitationRepository.ListByProject")
}
func (m *ProjectInvitationRepository) ListByStudent(ctx context.Context, studentID string) ([]*port.StudentInvitationView, error) {
	if m.ListByStudentFn != nil {
		return m.ListByStudentFn(ctx, studentID)
	}
	return nil, unimplemented("ProjectInvitationRepository.ListByStudent")
}
func (m *ProjectInvitationRepository) UpdateStatus(ctx context.Context, id string, status constants.InvitationStatus) (*entities.ProjectInvitation, error) {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, status)
	}
	return nil, unimplemented("ProjectInvitationRepository.UpdateStatus")
}
func (m *ProjectInvitationRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return unimplemented("ProjectInvitationRepository.Delete")
}

// ProfessorRepository is a mock of port.ProfessorRepository.
type ProfessorRepository struct {
	CreateFn      func(ctx context.Context, p *entities.Professor) (*entities.Professor, error)
	GetByIDFn     func(ctx context.Context, id string) (*entities.Professor, error)
	GetByUserIDFn func(ctx context.Context, userID string) (*entities.Professor, error)
	UpdateFn      func(ctx context.Context, p *entities.Professor) (*entities.Professor, error)
}

func (m *ProfessorRepository) Create(ctx context.Context, p *entities.Professor) (*entities.Professor, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, p)
	}
	return nil, unimplemented("ProfessorRepository.Create")
}
func (m *ProfessorRepository) GetByID(ctx context.Context, id string) (*entities.Professor, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, unimplemented("ProfessorRepository.GetByID")
}
func (m *ProfessorRepository) GetByUserID(ctx context.Context, userID string) (*entities.Professor, error) {
	if m.GetByUserIDFn != nil {
		return m.GetByUserIDFn(ctx, userID)
	}
	return nil, unimplemented("ProfessorRepository.GetByUserID")
}
func (m *ProfessorRepository) Update(ctx context.Context, p *entities.Professor) (*entities.Professor, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, p)
	}
	return nil, unimplemented("ProfessorRepository.Update")
}

// UserRepository is a mock of port.UserRepository.
type UserRepository struct {
	GetByIDFn    func(ctx context.Context, id string) (*entities.User, error)
	GetByEmailFn func(ctx context.Context, email string) (*entities.User, error)
	CreateFn     func(ctx context.Context, u *entities.User) (*entities.User, error)
	UpdateFn     func(ctx context.Context, u *entities.User) (*entities.User, error)
	DeleteFn     func(ctx context.Context, id string) error
}

func (m *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, unimplemented("UserRepository.GetByID")
}
func (m *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return nil, unimplemented("UserRepository.GetByEmail")
}
func (m *UserRepository) Create(ctx context.Context, u *entities.User) (*entities.User, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, u)
	}
	return nil, unimplemented("UserRepository.Create")
}
func (m *UserRepository) Update(ctx context.Context, u *entities.User) (*entities.User, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, u)
	}
	return nil, unimplemented("UserRepository.Update")
}
func (m *UserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return unimplemented("UserRepository.Delete")
}

// PasswordHasher is a mock of port.PasswordHasher.
type PasswordHasher struct {
	HashFn   func(plaintext string) (string, error)
	VerifyFn func(hashed, plaintext string) error
}

func (m *PasswordHasher) Hash(plaintext string) (string, error) {
	if m.HashFn != nil {
		return m.HashFn(plaintext)
	}
	return "", unimplemented("PasswordHasher.Hash")
}
func (m *PasswordHasher) Verify(hashed, plaintext string) error {
	if m.VerifyFn != nil {
		return m.VerifyFn(hashed, plaintext)
	}
	return unimplemented("PasswordHasher.Verify")
}

// TokenService is a mock of port.TokenService.
type TokenService struct {
	GenerateAccessFn  func(userID, email, role string) (string, error)
	GenerateRefreshFn func(userID string) (string, error)
	VerifyAccessFn    func(token string) (*port.AccessClaims, error)
	VerifyRefreshFn   func(token string) (*port.RefreshClaims, error)
}

func (m *TokenService) GenerateAccess(userID, email, role string) (string, error) {
	if m.GenerateAccessFn != nil {
		return m.GenerateAccessFn(userID, email, role)
	}
	return "", unimplemented("TokenService.GenerateAccess")
}
func (m *TokenService) GenerateRefresh(userID string) (string, error) {
	if m.GenerateRefreshFn != nil {
		return m.GenerateRefreshFn(userID)
	}
	return "", unimplemented("TokenService.GenerateRefresh")
}
func (m *TokenService) VerifyAccess(token string) (*port.AccessClaims, error) {
	if m.VerifyAccessFn != nil {
		return m.VerifyAccessFn(token)
	}
	return nil, unimplemented("TokenService.VerifyAccess")
}
func (m *TokenService) VerifyRefresh(token string) (*port.RefreshClaims, error) {
	if m.VerifyRefreshFn != nil {
		return m.VerifyRefreshFn(token)
	}
	return nil, unimplemented("TokenService.VerifyRefresh")
}

// RegisterCoreService is a mock of port.Service[RegisterInput, entities.User]
// used by the role-specific register usecases that delegate user creation to
// the core register service.
type RegisterCoreService struct {
	ExecuteFn func(ctx context.Context, in *inputoutput.RegisterInput) (*entities.User, error)
}

func (m *RegisterCoreService) Execute(ctx context.Context, in *inputoutput.RegisterInput) (*entities.User, error) {
	if m.ExecuteFn != nil {
		return m.ExecuteFn(ctx, in)
	}
	return nil, unimplemented("RegisterCoreService.Execute")
}

// StudentsRepository is a mock of port.StudentsRepository.
type StudentsRepository struct {
	CreateStudentFn      func(ctx context.Context, s *entities.Student) (*entities.Student, error)
	GetStudentByIDFn     func(ctx context.Context, id string) (*entities.Student, error)
	GetStudentByUserIDFn func(ctx context.Context, userID string) (*entities.Student, error)
	ListStudentsFn       func(ctx context.Context) ([]*entities.Student, error)
	UpdateStudentFn      func(ctx context.Context, s *entities.Student) (*entities.Student, error)
}

func (m *StudentsRepository) CreateStudent(ctx context.Context, s *entities.Student) (*entities.Student, error) {
	if m.CreateStudentFn != nil {
		return m.CreateStudentFn(ctx, s)
	}
	return nil, unimplemented("StudentsRepository.CreateStudent")
}
func (m *StudentsRepository) GetStudentByID(ctx context.Context, id string) (*entities.Student, error) {
	if m.GetStudentByIDFn != nil {
		return m.GetStudentByIDFn(ctx, id)
	}
	return nil, unimplemented("StudentsRepository.GetStudentByID")
}
func (m *StudentsRepository) GetStudentByUserID(ctx context.Context, userID string) (*entities.Student, error) {
	if m.GetStudentByUserIDFn != nil {
		return m.GetStudentByUserIDFn(ctx, userID)
	}
	return nil, unimplemented("StudentsRepository.GetStudentByUserID")
}
func (m *StudentsRepository) ListStudents(ctx context.Context) ([]*entities.Student, error) {
	if m.ListStudentsFn != nil {
		return m.ListStudentsFn(ctx)
	}
	return nil, unimplemented("StudentsRepository.ListStudents")
}
func (m *StudentsRepository) UpdateStudent(ctx context.Context, s *entities.Student) (*entities.Student, error) {
	if m.UpdateStudentFn != nil {
		return m.UpdateStudentFn(ctx, s)
	}
	return nil, unimplemented("StudentsRepository.UpdateStudent")
}

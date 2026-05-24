package inputoutput

// --- Create project ---

type CreateProjectInput struct {
	// ProfessorUserID is sourced from the JWT, not the request body.
	ProfessorUserID string `json:"-"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Slots           int    `json:"slots"`
}

type ProjectOutput struct {
	ID          string                 `json:"id"`
	ProfessorID string                 `json:"professor_id"`
	Professor   *ProjectProfessorBrief `json:"professor"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Slots       int                    `json:"slots"`
	Status      string                 `json:"status"`
}

type ProjectProfessorBrief struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	University string `json:"university"`
	Department string `json:"department"`
}

// --- Listing ---

type ListProjectsInput struct{}

type ListProjectsOutput struct {
	Projects []ProjectOutput `json:"projects"`
}

// --- Get one project ---

type GetProjectInput struct {
	ProjectID string `json:"-"`
}

// --- Apply to project (student) ---

type ApplyProjectInput struct {
	// StudentUserID is sourced from the JWT.
	StudentUserID string `json:"-"`
	ProjectID     string `json:"-"`
	Message       string `json:"message"`
}

type ProjectApplicationOutput struct {
	ID        string                `json:"id"`
	Student   *ProjectStudentBrief  `json:"student,omitempty"`
	Status    string                `json:"status"`
	Message   string                `json:"message,omitempty"`
	Project   *ProjectShortForApp   `json:"project,omitempty"`
}

type ProjectStudentBrief struct {
	StudentID         string `json:"student_id"`
	UserID            string `json:"user_id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	University        string `json:"university"`
	Department        string `json:"department"`
	ResearchInterests string `json:"research_interests"`
}

type ProjectShortForApp struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// --- List applications for a project (professor) ---

type ListApplicationsByProjectInput struct {
	// ProfessorUserID is the caller, used to verify project ownership.
	ProfessorUserID string `json:"-"`
	ProjectID       string `json:"-"`
}

type ListApplicationsOutput struct {
	Applications []*ProjectApplicationOutput `json:"applications"`
}

type ListProjectApplicationsByProjectOutput struct {
	ApprovedApplications []*ProjectApplicationOutput `json:"approved_applications"`
	PendingApplications []*ProjectApplicationOutput `json:"pending_applications"`
} 

// --- Review an application (professor approves/rejects) ---

type ReviewApplicationInput struct {
	ProfessorUserID string `json:"-"`
	ProjectID       string `json:"-"`
	ApplicationID   string `json:"-"`
	Status          string `json:"status"`
}

// --- Withdraw (student) ---

type WithdrawApplicationInput struct {
	StudentUserID string `json:"-"`
	ProjectID     string `json:"-"`
	ApplicationID string `json:"-"`
}

type EmptyOutput struct {
	Message string `json:"message"`
}

// --- Student lists their own applications ---

type ListMyApplicationsInput struct {
	StudentUserID string `json:"-"`
}

type CheckApplicationStatusInput struct {
	StudentUserID string `json:"-"`
	ProjectID     string `json:"-"`
}

type CheckApplicationStatusOutput struct {
	Exists bool `json:"exists"`
}

// --- Project members ---

type ProjectMemberOutput struct {
	ID          string               `json:"id"`
	Student     *ProjectStudentBrief `json:"student,omitempty"`
	Source      string               `json:"source"`
	SourceRefID string               `json:"source_ref_id,omitempty"`
	Status      string               `json:"status"`
	JoinedAt    string               `json:"joined_at,omitempty"`
	LeftAt      string               `json:"left_at,omitempty"`
	Project     *ProjectShortForApp  `json:"project,omitempty"`
}

type ListMembersByProjectInput struct {
	ProjectID string `json:"-"`
}

type ListMembersOutput struct {
	Members []*ProjectMemberOutput `json:"members"`
}

type RemoveMemberInput struct {
	ProfessorUserID string `json:"-"`
	ProjectID       string `json:"-"`
	MemberID        string `json:"-"`
}

type LeaveProjectInput struct {
	StudentUserID string `json:"-"`
	ProjectID     string `json:"-"`
}

type ListMyProjectsInput struct {
	StudentUserID string `json:"-"`
}

type ListMyProjectsOutput struct {
	Memberships []*ProjectMemberOutput `json:"memberships"`
}

// --- Project invitations ---

type SendInvitationInput struct {
	ProfessorUserID string `json:"-"`
	ProjectID       string `json:"-"`
	StudentID       string `json:"student_id"`
	Message         string `json:"message"`
}

type ProjectInvitationOutput struct {
	ID          string               `json:"id"`
	Student     *ProjectStudentBrief `json:"student,omitempty"`
	Project     *ProjectShortForApp  `json:"project,omitempty"`
	Status      string               `json:"status"`
	Message     string               `json:"message,omitempty"`
	RespondedAt string               `json:"responded_at,omitempty"`
}

type ListInvitationsByProjectInput struct {
	ProfessorUserID string `json:"-"`
	ProjectID       string `json:"-"`
}

type ListInvitationsOutput struct {
	Invitations []*ProjectInvitationOutput `json:"invitations"`
}

type CancelInvitationInput struct {
	ProfessorUserID string `json:"-"`
	ProjectID       string `json:"-"`
	InvitationID    string `json:"-"`
}

type ListMyInvitationsInput struct {
	StudentUserID string `json:"-"`
}

type RespondInvitationInput struct {
	StudentUserID string `json:"-"`
	InvitationID  string `json:"-"`
	Status        string `json:"status"`
}

// --- Directory / lookup ---

type ListStudentsInput struct{}

type ListStudentsOutput struct {
	Students []*ProjectStudentBrief `json:"students"`
}
package inputoutput

// RegisterInput defines the input for user registration
type RegisterInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterProfessorInput struct {
	RegisterInput
	University string `json:"university"`
	Department string `json:"department"`
}

type RegisterStudentInput struct {
	RegisterInput
	University        string `json:"university"`
	Department        string `json:"department"`
	ResearchInterests string `json:"research_interests"`
}

// RegisterOutput defines the output for user registration
type RegisterOutput struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

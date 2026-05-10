package inputoutput

// ProfileInput is the request payload for fetching a profile.
// The user id is sourced from the decoded JWT, not the request body.
type ProfileInput struct {
	UserID string
}

// ProfessorProfileOutput is the response payload for the professor profile endpoint.
type ProfessorProfileOutput struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	University string `json:"university"`
	Department string `json:"department"`
}

// StudentProfileOutput is the response payload for the student profile endpoint.
type StudentProfileOutput struct {
	ID                string `json:"id"`
	UserID            string `json:"user_id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Role              string `json:"role"`
	University        string `json:"university"`
	Department        string `json:"department"`
	ResearchInterests string `json:"research_interests"`
}

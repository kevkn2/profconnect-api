package inputoutput

// LoginInput represents the login request payload
type LoginInput struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// LoginOutput represents the login response payload
type LoginOutput struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Type  string `json:"type" example:"Bearer"`
}

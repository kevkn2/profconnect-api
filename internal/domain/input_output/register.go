package inputoutput

// RegisterInput defines the input for user registration
type RegisterInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// RegisterOutput defines the output for user registration
type RegisterOutput struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

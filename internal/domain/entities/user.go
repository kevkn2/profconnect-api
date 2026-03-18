package entities

// User represents a user entity in the system
type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"` // Never expose in JSON
}

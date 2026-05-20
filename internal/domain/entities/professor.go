package entities

type Professor struct {
	ID         string `json:"id"`
	User       *User  `json:"user"`
	University string `json:"university"`
	Department string `json:"department"`
}
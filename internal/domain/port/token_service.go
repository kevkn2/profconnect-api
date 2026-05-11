package port

type AccessClaims struct {
	UserID string
	Email  string
	Role   string
}

type RefreshClaims struct {
	UserID string
}

type TokenService interface {
	GenerateAccess(userID, email, role string) (string, error)
	GenerateRefresh(userID string) (string, error)
	VerifyAccess(token string) (*AccessClaims, error)
	VerifyRefresh(token string) (*RefreshClaims, error)
}

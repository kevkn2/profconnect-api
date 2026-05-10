package crypto

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTSecret        = "your-secret-key-change-in-production" // TODO: Load from environment variable
	TokenTTL         = 24 * time.Hour
	RefreshTokenTTL  = 7 * 24 * time.Hour
	refreshTokenType = "refresh"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshClaims is the minimal claim set carried by a refresh token.
// It deliberately omits email/role — those are re-fetched from the
// database when a new access token is issued, so they can't drift.
type RefreshClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token for the given user ID, email, and role
func GenerateJWT(userID, email, role string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "profconnect-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT verifies and parses a JWT token, returning the custom claims
func VerifyJWT(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateRefreshToken generates a refresh token carrying only the user id.
// Refresh tokens have a longer TTL than access tokens.
func GenerateRefreshToken(userID string) (string, error) {
	claims := RefreshClaims{
		UserID:    userID,
		TokenType: refreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "profconnect-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// VerifyRefreshToken verifies and parses a refresh token, returning the claims.
// It additionally checks that the token_type marker matches "refresh" so that
// an access token cannot be replayed against the refresh endpoint.
func VerifyRefreshToken(tokenString string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != refreshTokenType {
		return nil, errors.New("not a refresh token")
	}

	return claims, nil
}

// GetTokenTTL returns the access token time-to-live duration
func GetTokenTTL() time.Duration {
	return TokenTTL
}

// GetRefreshTokenTTL returns the refresh token time-to-live duration
func GetRefreshTokenTTL() time.Duration {
	return RefreshTokenTTL
}

// SetJWTSecret allows setting the JWT secret from environment or config
func SetJWTSecret(secret string) {
	if secret != "" {
		JWTSecret = secret
	}
}

package token

import (
	"errors"
	"time"

	"profconnect-api/internal/domain/port"

	"github.com/golang-jwt/jwt/v5"
)

const refreshTokenType = "refresh"

type accessClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

func NewJWTService(secret string, accessTTL, refreshTTL time.Duration, issuer string) port.TokenService {
	return &jwtService{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		issuer:          issuer,
	}
}

func (j *jwtService) GenerateAccess(userID, email, role string) (string, error) {
	claims := accessClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *jwtService) GenerateRefresh(userID string) (string, error) {
	claims := refreshClaims{
		UserID:    userID,
		TokenType: refreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *jwtService) VerifyAccess(tokenString string) (*port.AccessClaims, error) {
	claims := &accessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, j.keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &port.AccessClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

func (j *jwtService) VerifyRefresh(tokenString string) (*port.RefreshClaims, error) {
	claims := &refreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, j.keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != refreshTokenType {
		return nil, errors.New("not a refresh token")
	}
	return &port.RefreshClaims{UserID: claims.UserID}, nil
}

func (j *jwtService) keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return j.secret, nil
}

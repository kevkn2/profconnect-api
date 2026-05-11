package middleware

import (
	"strings"

	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/port"

	"github.com/gofiber/fiber/v3"
)

// Locals keys used to share decoded JWT data with handlers.
const (
	LocalsUserID = "user_id"
	LocalsEmail  = "email"
	LocalsRole   = "role"
)

type errorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

// Auth groups authentication middlewares bound to a TokenService.
type Auth struct {
	tokenService port.TokenService
}

// NewAuth constructs an Auth middleware bundle.
func NewAuth(tokenService port.TokenService) *Auth {
	return &Auth{tokenService: tokenService}
}

// JWT decodes the bearer token, verifies its signature, and stores
// the user_id, email, and role in fiber Locals for downstream handlers.
func (a *Auth) JWT() fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{
				Message: "missing authorization header",
				Type:    "unauthorized",
			})
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{
				Message: "invalid authorization header format",
				Type:    "unauthorized",
			})
		}

		claims, err := a.tokenService.VerifyAccess(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{
				Message: "invalid or expired token",
				Type:    "unauthorized",
			})
		}

		c.Locals(LocalsUserID, claims.UserID)
		c.Locals(LocalsEmail, claims.Email)
		c.Locals(LocalsRole, claims.Role)

		return c.Next()
	}
}

// RequireRole returns a middleware that enforces the caller's role matches
// the expected role. Must be chained after JWT.
func RequireRole(expected constants.Roles) fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals(LocalsRole).(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{
				Message: "missing role in token",
				Type:    "unauthorized",
			})
		}

		if role != string(expected) {
			return c.Status(fiber.StatusForbidden).JSON(errorResponse{
				Message: "insufficient permissions for this resource",
				Type:    "forbidden",
			})
		}

		return c.Next()
	}
}

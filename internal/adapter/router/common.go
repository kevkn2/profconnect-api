package router

import (
	authHandler "profconnect-api/internal/adapter/handler/auth"
	professorHandler "profconnect-api/internal/adapter/handler/professor"
	studentHandler "profconnect-api/internal/adapter/handler/student"
	"profconnect-api/internal/adapter/middleware"

	"github.com/gofiber/fiber/v3"
)

// Router handles all route registrations
type Router struct {
	app       *fiber.App
	auth      *authHandler.Handler
	professor *professorHandler.Handler
	student   *studentHandler.Handler
	authMW    *middleware.Auth
}

// NewRouter creates a new router instance
func NewRouter(
	app *fiber.App,
	auth *authHandler.Handler,
	professor *professorHandler.Handler,
	student *studentHandler.Handler,
	authMW *middleware.Auth,
) *Router {
	return &Router{
		app:       app,
		auth:      auth,
		professor: professor,
		student:   student,
		authMW:    authMW,
	}
}

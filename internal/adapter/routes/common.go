package routes

import (
	"profconnect-api/internal/adapter/handler"

	"github.com/gofiber/fiber/v3"
)

// Router handles all route registrations
type Router struct {
	app     *fiber.App
	handler *handler.Handler
}

// NewRouter creates a new router instance
func NewRouter(app *fiber.App, handler *handler.Handler) *Router {
	return &Router{
		app:     app,
		handler: handler,
	}
}

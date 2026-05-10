package routes

import "profconnect-api/internal/adapter/handler"

func (r *Router) RegisterSwaggerRoutes() {
	r.app.Get("/swagger", handler.Swagger)
}

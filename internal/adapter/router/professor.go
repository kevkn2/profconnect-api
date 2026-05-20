package router

import (
	"profconnect-api/internal/adapter/middleware"
	"profconnect-api/internal/domain/constants"
)

func (r *Router) RegisterProfessorRoutes() {
	professorGroup := r.app.Group("/api/professor", r.authMW.JWT(), middleware.RequireRole(constants.Professor))

	professorGroup.Get("/profile", r.professor.Profile)
	professorGroup.Post("/projects", r.professor.CreateProject)
	professorGroup.Get("/projects/:id/applications", r.professor.ListApplicationsByProject)
	professorGroup.Patch("/projects/:id/applications/:appId", r.professor.ReviewApplication)
}
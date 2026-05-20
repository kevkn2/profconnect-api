package router

import (
	"profconnect-api/internal/adapter/middleware"
	"profconnect-api/internal/domain/constants"
)

func (r *Router) RegisterStudentRoutes() {
	studentGroup := r.app.Group("/api/student", r.authMW.JWT(), middleware.RequireRole(constants.Student))

	studentGroup.Get("/profile", r.student.Profile)
	studentGroup.Get("/applications", r.student.ListMyApplications)
	studentGroup.Post("/projects/:id/applications", r.student.ApplyProject)
	studentGroup.Delete("/projects/:id/applications/:appId", r.student.WithdrawApplication)
}
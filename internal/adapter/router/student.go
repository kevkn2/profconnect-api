package router

import (
	"profconnect-api/internal/adapter/middleware"
	"profconnect-api/internal/domain/constants"
)

func (r *Router) RegisterStudentRoutes() {
	studentGroup := r.app.Group("/api/student", r.authMW.JWT(), middleware.RequireRole(constants.Student))

	studentGroup.Get("/profile", r.student.Profile)
	studentGroup.Get("/applications", r.student.ListMyApplications)
	studentGroup.Get("/projects/:id/applications", r.student.ListApplicationsPerProject)
	studentGroup.Post("/projects/:id/applications", r.student.ApplyProject)
	studentGroup.Delete("/projects/:id/applications/:appId", r.student.WithdrawApplication)
	studentGroup.Get("/projects", r.student.ListMyProjects)
	studentGroup.Delete("/projects/:id/membership", r.student.LeaveProject)
	studentGroup.Get("/invitations", r.student.ListMyInvitations)
	studentGroup.Patch("/invitations/:invId", r.student.RespondInvitation)
}

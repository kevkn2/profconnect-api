package routes

import (
	"profconnect-api/internal/adapter/handler"
	"profconnect-api/internal/adapter/middleware"
	"profconnect-api/internal/domain/constants"
)

func (r *Router) RegisterGeneralRoutes() {
	r.app.Get("/", handler.HelloWorld)
}

func (r *Router) RegisterAuthRoutes() {
	authGroup := r.app.Group("/api/auth")

	authGroup.Post("/register/admin", r.auth.RegisterAdmin)
	authGroup.Post("/register/professor", r.auth.RegisterProfessor)
	authGroup.Post("/register/student", r.auth.RegisterStudent)
	authGroup.Post("/login", r.auth.Login)
	authGroup.Post("/refresh", r.auth.Refresh)
}

func (r *Router) RegisterStudentRoutes() {
	studentGroup := r.app.Group("/api/student", middleware.JWTAuth(), middleware.RequireRole(constants.Student))

	studentGroup.Get("/profile", r.student.Profile)
}

func (r *Router) RegisterProfessorRoutes() {
	professorGroup := r.app.Group("/api/professor", middleware.JWTAuth(), middleware.RequireRole(constants.Professor))

	professorGroup.Get("/profile", r.professor.Profile)
}

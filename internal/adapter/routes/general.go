package routes

import (
	"profconnect-api/internal/adapter/middleware"
	"profconnect-api/internal/domain/constants"
)

func (r *Router) RegisterGeneralRoutes() {
	r.app.Get("/", r.handler.HelloWorld)

	r.app.Post("/api/register/admin", r.handler.RegisterAdmin)
	r.app.Post("/api/register/professor", r.handler.RegisterProfessor)
	r.app.Post("/api/register/student", r.handler.RegisterStudent)
	r.app.Post("/api/login", r.handler.Login)

	r.app.Get(
		"/api/profile/professor",
		r.handler.ProfessorProfile,
		middleware.JWTAuth(),
		middleware.RequireRole(constants.Professor),
	)
	r.app.Get(
		"/api/profile/student",
		r.handler.StudentProfile,
		middleware.JWTAuth(),
		middleware.RequireRole(constants.Student),
	)
}

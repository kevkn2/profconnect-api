package router

func (r *Router) RegisterAuthRoutes() {
	authGroup := r.app.Group("/api/auth")

	authGroup.Post("/register/admin", r.auth.RegisterAdmin)
	authGroup.Post("/register/professor", r.auth.RegisterProfessor)
	authGroup.Post("/register/student", r.auth.RegisterStudent)
	authGroup.Post("/login", r.auth.Login)
	authGroup.Post("/refresh", r.auth.Refresh)
}
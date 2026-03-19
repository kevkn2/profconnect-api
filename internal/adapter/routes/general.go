package routes

func (r *Router) RegisterGeneralRoutes() {
	r.app.Get("/", r.handler.HelloWorld)
	r.app.Post("/api/register/admin", r.handler.RegisterAdmin)
	r.app.Post("/api/register/professor", r.handler.RegisterProfessor)
	r.app.Post("/api/login", r.handler.Login)
}

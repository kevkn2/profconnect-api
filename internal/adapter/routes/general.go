package routes

func (r *Router) RegisterGeneralRoutes() {
	r.app.Get("/", r.handler.HelloWorld)
	r.app.Post("/api/register", r.handler.Register)
	r.app.Post("/api/login", r.handler.Login)
}

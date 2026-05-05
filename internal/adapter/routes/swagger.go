package routes

func (r *Router) RegisterSwaggerRoutes() {
	r.app.Get("/swagger", r.handler.Swagger) // default
}

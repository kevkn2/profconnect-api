package router

func (r *Router) RegisterProjectRoutes() {
	// Browsing and viewing projects is open to any authenticated user.
	authed := r.app.Group("/api/projects", r.authMW.JWT())
	authed.Get("/", r.project.ListProjects)
	authed.Get("/:id", r.project.GetProject)
}
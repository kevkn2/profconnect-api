package router

import (
	"profconnect-api/internal/adapter/handler"
)

func (r *Router) RegisterGeneralRoutes() {
	r.app.Get("/", handler.HelloWorld)
}

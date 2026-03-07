package routes

import (
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/http/controllers"
)

// RegisterAPI defines API routes under the /api prefix.
func RegisterAPI(r *router.Router) {
	api := r.Group("/api")
	api.APIResource("/posts", &controllers.PostController{})
}

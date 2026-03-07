package routes

import (
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/http/controllers"
)

// RegisterWeb defines web (HTML) routes.
func RegisterWeb(r *router.Router) {
	r.Get("/", controllers.Home)
}

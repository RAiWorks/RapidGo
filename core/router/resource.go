package router

import "github.com/gin-gonic/gin"

// ResourceController defines the interface for RESTful controllers.
// Implement all 7 methods for full CRUD with form routes,
// or use APIResource to skip Create/Edit form routes.
type ResourceController interface {
	Index(c *gin.Context)   // GET    /resource
	Create(c *gin.Context)  // GET    /resource/create  (SSR form)
	Store(c *gin.Context)   // POST   /resource
	Show(c *gin.Context)    // GET    /resource/:id
	Edit(c *gin.Context)    // GET    /resource/:id/edit (SSR form)
	Update(c *gin.Context)  // PUT    /resource/:id
	Destroy(c *gin.Context) // DELETE /resource/:id
}

// Resource registers all 7 RESTful routes for a controller on the router.
func (r *Router) Resource(path string, ctrl ResourceController) {
	r.engine.GET(path, ctrl.Index)
	r.engine.GET(path+"/create", ctrl.Create)
	r.engine.POST(path, ctrl.Store)
	r.engine.GET(path+"/:id", ctrl.Show)
	r.engine.GET(path+"/:id/edit", ctrl.Edit)
	r.engine.PUT(path+"/:id", ctrl.Update)
	r.engine.DELETE(path+"/:id", ctrl.Destroy)
}

// APIResource registers 5 RESTful routes (no Create/Edit form routes) on the router.
func (r *Router) APIResource(path string, ctrl ResourceController) {
	r.engine.GET(path, ctrl.Index)
	r.engine.POST(path, ctrl.Store)
	r.engine.GET(path+"/:id", ctrl.Show)
	r.engine.PUT(path+"/:id", ctrl.Update)
	r.engine.DELETE(path+"/:id", ctrl.Destroy)
}

// Resource registers all 7 RESTful routes for a controller on the group.
func (g *RouteGroup) Resource(path string, ctrl ResourceController) {
	g.group.GET(path, ctrl.Index)
	g.group.GET(path+"/create", ctrl.Create)
	g.group.POST(path, ctrl.Store)
	g.group.GET(path+"/:id", ctrl.Show)
	g.group.GET(path+"/:id/edit", ctrl.Edit)
	g.group.PUT(path+"/:id", ctrl.Update)
	g.group.DELETE(path+"/:id", ctrl.Destroy)
}

// APIResource registers 5 RESTful routes (no Create/Edit form routes) on the group.
func (g *RouteGroup) APIResource(path string, ctrl ResourceController) {
	g.group.GET(path, ctrl.Index)
	g.group.POST(path, ctrl.Store)
	g.group.GET(path+"/:id", ctrl.Show)
	g.group.PUT(path+"/:id", ctrl.Update)
	g.group.DELETE(path+"/:id", ctrl.Destroy)
}

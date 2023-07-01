package notion

import "github.com/gin-gonic/gin"

func (a NotionApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("notion")
	{
		App.Routes.GET("/import", App.RouteHandler.ImportNotion)
	}
}

package parser2

import "github.com/gin-gonic/gin"

func (a Parser2App) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("parser")
	{
		App.Routes.GET("/markdown/:canvasBranchID", App.RouteHandler.GetCanvasBranchMarkdownFile)
		App.Routes.POST("/import-notion", App.RouteHandler.ImportNotion)
		App.Routes.POST("/import-file", App.RouteHandler.ImportFile)
		App.Routes.GET("/export-studio", App.RouteHandler.ExportStudio)
	}
}

package permissions

import "github.com/gin-gonic/gin"

func (a permissionApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("permission")
	{
		App.Routes.GET("/studio", App.RouteHandler.getStudioPermission)
		App.Routes.GET("/collection", App.RouteHandler.getCollectionPermission)
		App.Routes.GET("/canvas/:collectionId", App.RouteHandler.getCanvasPermission)
		App.Routes.GET("/canvas/:collectionId/:parentCanvasId", App.RouteHandler.getSubCanvasPermission)
		App.Routes.GET("/invalidate-cache", App.RouteHandler.invalidateCache)
		App.Routes.GET("/flush-cache", App.RouteHandler.flushCache)
	}
}

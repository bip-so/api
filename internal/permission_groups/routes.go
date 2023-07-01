package permissiongroup

import (
	"github.com/gin-gonic/gin"
)

func RegisterPermsGroupRoutes(router *gin.RouterGroup) {
	gpr := router.Group("permissions-schema")
	{
		// Anonymous Calls
		gpr.GET("/studio/schema", PermissionGroupRouteHandler.StudioPermsSchema)
		gpr.GET("/collection/schema", PermissionGroupRouteHandler.CollectionPermsSchema)
		gpr.GET("/canvasBranch/schema", PermissionGroupRouteHandler.CanvasBranchPermsSchema)
	}
}

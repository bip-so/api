package studiopermissions

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func RegisterStudioPermissionRoutes(router *gin.RouterGroup) {
	studioPerm := router.Group("studiopermission")
	{
		// no specific check
		studioPerm.GET("/getAll", App.RouteHandler.getStudioPermissionsRouteHandler)
		// STUDIO_MANAGE_PERMS
		studioPerm.POST("/update", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.updateStudioPermissionsRouteHandler)
		// STUDIO_MANAGE_PERMS
		studioPerm.DELETE("/:studioPermissionID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.deleteStudioPermissionsRouteHandler)

	}
}

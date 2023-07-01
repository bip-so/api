package canvasbranchpermissions

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func RegisterCanvasBranchPermissionRoutes(router *gin.RouterGroup) {
	canvasBranchPerm := router.Group("canvasbranchpermission")
	{
		// CANVAS_BRANCH_MANAGE_PERMS, cannot change BipAdmin role perm -> Role is not updated here
		canvasBranchPerm.POST("/update", middlewares.TokenAuthorizationMiddleware(), CanvasBranchPermissionsRouteHandler.createCanvasBranchPermissionsRouteHandler)
		canvasBranchPerm.POST("/bulk-update", middlewares.TokenAuthorizationMiddleware(), CanvasBranchPermissionsRouteHandler.BulkCreateCanvasBranchPermissionsRouteHandler)
		// CANVAS_BRANCH_VIEW
		canvasBranchPerm.GET("/:canvasBranchId", CanvasBranchPermissionsRouteHandler.getCanvasBranchPermissionsRouteHandler)
		// CANVAS_BRANCH_MANAGE_PERMS, cannot remove BipAdmin role
		canvasBranchPerm.DELETE("/:canvasBranchPermissionId", middlewares.TokenAuthorizationMiddleware(), CanvasBranchPermissionsRouteHandler.deleteCanvasBranchPermissionsRouteHandler)
		// CANVAS_BRANCH_MANAGE_PERMS
		canvasBranchPerm.POST("/inherit/:canvasBranchId", CanvasBranchPermissionsRouteHandler.inheritParentPermissionsRouteHandler)
	}
}

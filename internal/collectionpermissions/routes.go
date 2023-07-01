package collectionpermissions

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func RegisterCollectionPermissionRoutes(router *gin.RouterGroup) {
	collectionPerm := router.Group("collectionpermission")
	{
		// COLLECTION_MANAGE_PERMS, cannot update or delete bipAdmin role
		collectionPerm.POST("/update", middlewares.TokenAuthorizationMiddleware(), CollectionPermissionsRouteHandler.createCollectionPermissionsRouteHandler)
		// COLLECTION_VIEW_METADATA
		collectionPerm.GET("/:collectionId", CollectionPermissionsRouteHandler.getCollectionPermissionsRouteHandler)
		// COLLECTION_MANAGE_PERMS, cannot update or delete bipAdmin role
		collectionPerm.DELETE("/:collectionPermissionId", middlewares.TokenAuthorizationMiddleware(), CollectionPermissionsRouteHandler.deleteCollectionPermissionsRouteHandler)

	}
}

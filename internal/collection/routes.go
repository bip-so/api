package collection

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a CollectionApp) RegisterCollectionRoutes(router *gin.RouterGroup) {
	App.Routes = *router.Group("collection")
	{
		App.Routes.GET("/test", App.RouteHandler.test)

		// STUDIO_CREATE_COLLECTION
		App.Routes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.createCollectionRoute)
		// No specific check
		App.Routes.GET("/get", App.RouteHandler.getCollectionRoute)
		// COLLECTION_EDIT_NAME
		App.Routes.PUT("/update", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.updateCollectionRoute)
		// COLLECTION_MANAGE_PERMS
		App.Routes.POST("/:collectionId/visibility", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.manageVisibilityCollectionRoute)
		// COLLECTION_DELETE
		App.Routes.DELETE("/delete/:collectionId", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.deleteCollectionRoute)
		// STUDIO_CHANGE_CANVAS_COLLECTION_POSITION
		App.Routes.POST("/move", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.moveCollectionRoute)
		//App.Routes.POST("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.moveCollectionRoute)
		// STUDIO_ADD_REMOVE_USER_TO_ROLE
		App.Routes.GET("/user/:userId", App.RouteHandler.getStudioMemberCollections)
		App.Routes.GET("/role/:roleId", App.RouteHandler.getStudioRoleCollections)
		App.Routes.GET("/next-prev/:collectionId", App.RouteHandler.getNextPrevCollection)
	}
}

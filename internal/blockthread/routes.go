package blockthread

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a blockThreadApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("block-thread")
	{
		// PermissionGroup : CANVAS_BRANCH_VIEW
		App.Routes.GET("/:blockThreadID", App.RouteHandler.Get)
		// PermissionGroup: CANVAS_BRANCH_VIEW
		App.Routes.GET("/branch/:canvasBranchID", App.RouteHandler.GetByBranch)
		// Todo: Resolve
		// Creator of comment
		// CANVAS_BRANCH_MANAGE_CONTENT
		App.Routes.POST("/:blockThreadID/resolve", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Resolve)
		// CANVAS_BRANCH_ADD_COMMENT
		App.Routes.POST("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Create)
		// creator of comment
		App.Routes.PATCH("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Update)
		// Creator of comment
		// CANVAS_BRANCH_MANAGE_CONTENT
		App.Routes.DELETE("/:blockThreadID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Delete)

		//App.Routes.POST("/:blockThreadID/", App.RouteHandler.Get)

	}
}

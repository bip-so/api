package blockThreadCommentcomment

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a blockThreadCommentApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("block-thread-comment")
	{
		// CANVAS_BRANCH_VIEW
		App.Routes.GET("/:blockThreadID", App.RouteHandler.Get)
		App.Routes.GET("reply/:blockThreadID/:parentCommentID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.GetReply) // not using this in FE remove later
		// CANVAS_BRANCH_ADD_COMMENT
		App.Routes.POST("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Create)
		// creator of comment
		App.Routes.PATCH("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Update)
		// Creator of comment
		// CANVAS_BRANCH_MANAGE_CONTENT
		App.Routes.DELETE("/:blockThreadCommentID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Delete)
	}
}

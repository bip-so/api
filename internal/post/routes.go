package post

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a postApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("posts")
	{
		// POST
		App.Routes.GET("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.GetAllPost)
		App.Routes.GET("/homepage", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.GetPostHomepage)
		App.Routes.GET("/:postID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.GetOnePost)
		App.Routes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.CreatePost)
		App.Routes.PATCH("/:postID/edit", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.UpdatePost)
		App.Routes.DELETE("/:postID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.DeletePost)

		App.Routes.POST("/:postID/add-reaction", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.AddReactionPost)
		App.Routes.POST("/:postID/remove-reaction", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.RemoveReactionPost)

		// Post Comment
		postCommentRoutes := App.Routes.Group("/:postID/comments")
		postCommentRoutes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.CreateCommentPost)
		postCommentRoutes.GET("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.GetAllPostComments)
		postCommentRoutes.DELETE("/:postCommentID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.DeleteCommentPost)
		postCommentRoutes.PATCH("/:postCommentID/edit", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.UpdateCommentPost)

		postCommentRoutes.POST("/:postCommentID/add-reaction", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.AddReactionPostComment)
		postCommentRoutes.POST("/:postCommentID/remove-reaction", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.RemoveReactionPostComment)
	}
}

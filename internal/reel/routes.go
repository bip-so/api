package reel

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a ReelApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("reels")
	{
		// Get all the Reels for this studio.
		// Studio is taken from Headers
		App.Routes.GET("/", App.RouteHandler.Get)
		App.Routes.GET("/:reelID", App.RouteHandler.GetOneReel)
		App.Routes.DELETE("/:reelID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.DeleteReel)
		// Todo: Add Pagination (50)
		App.Routes.GET("/popular", App.RouteHandler.GetPopularReels)
		// CANVAS_BRANCH_CREATE_REEL
		App.Routes.POST("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Create)
		// No specific check
		App.Routes.GET("/feed", App.RouteHandler.GetReelsFeed)

		comments := App.Routes.Group("/:reelID/comments")
		// CANVAS_BRANCH_COMMENT_ON_REEL
		comments.POST("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.NewCommentReel)
		// CANVAS_BRANCH_VIEW
		comments.GET("/", App.RouteHandler.GetCommentsReel)
		// Self reel comment creator can only delete it.
		comments.DELETE("/:reelCommentID", App.RouteHandler.DeleteReelComment)
	}
}

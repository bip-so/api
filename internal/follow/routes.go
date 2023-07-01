package follow

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a FollowApp) RegisterRoutes(router *gin.RouterGroup) {
	App.Routes = *router.Group("follow")
	{
		// users
		// Get Number of followers for a user & Get Number of follwed by this user]
		user := App.Routes.Group("user")
		{
			user.GET("/follow-count", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.followUserFollowCountRoute)
			user.POST("/follow", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.followUserRoute)
			// only self user can unfollow
			user.POST("/unfollow", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.unfollowUserRoute)
			user.GET("/list", App.RouteHandler.FollowList)
		}

		// studio
		studio := App.Routes.Group("studio")
		{
			studio.GET("/follower", App.RouteHandler.followStudioCountRoute)
			studio.POST("/follow", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.followUserStudioRoute)
			studio.POST("/unfollow", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.unfollowStudioRoute)
		}

		//
		//follow.POST("/user", App.RouteHandler.followUserRoute)
		//follow.POST("/studio", App.RouteHandler.followStudioRoute)
		//follow.POST("/topic", App.RouteHandler.followTopicRoute)
		//follow.GET("/studioFollowings", App.RouteHandler.studioFollowingsRoute)
		//follow.GET("/userFollowings", App.RouteHandler.userFollowingsRoute)
		//follow.GET("/getUserStudioFollowers", App.RouteHandler.userStudioFollowersRoute)
	}
}

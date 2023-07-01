package user

import (
	"github.com/gin-gonic/gin"
)

func (a UserApp) RegisterUser(router *gin.RouterGroup) {
	App.Routes = *router.Group("user")
	{
		App.Routes.POST("/setup", App.RouteHandler.setupUserRoute)
		App.Routes.GET("/info", App.RouteHandler.userInfoRoute)
		App.Routes.PUT("/update", App.RouteHandler.updateUserRoute)
		App.Routes.POST("/search", App.RouteHandler.UserSearchRoute)
		App.Routes.GET("/followers-list", App.RouteHandler.followerListRoute)
		App.Routes.GET("/following-list", App.RouteHandler.followingListRoute)
		//user.POST("/update-follow", App.RouteHandler.updateFollowerRoute)
		settings := App.Routes.Group("settings")
		settings.GET("", App.RouteHandler.GetUserSettingsRoute)
		settings.PATCH("", App.RouteHandler.UpdateUserSettingsRoute)
		contacts := App.Routes.Group("contacts")
		contacts.GET("", App.RouteHandler.GetUserContactsRoute)
		contacts.POST("", App.RouteHandler.UpdateUserContactsRoute)

	}
}

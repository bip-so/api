package global

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func RegisterGlobalRoutes(router *gin.RouterGroup) {
	studio := router.Group("studio")
	{
		// loggedin check
		studio.POST("/create", middlewares.TokenAuthorizationMiddleware(), GlobalRouteHandler.createStudioRouteHandler)

		//studio.PUT("/members/:userId", member.MemberRouteHandler.AddMemberToStudio)
	}

	user := router.Group("user")
	{
		user.GET("/popular", GlobalRouteHandler.popularUsersRoute)
	}
	// @todo: StudioFiles Model Future.

	global := router.Group("global")
	{
		global.GET("/check-handle", GlobalRouteHandler.checkHandleAvailableRouteHandler)
		global.GET("/search", GlobalRouteHandler.searchRouteHandler)
		global.POST("/addsearch", GlobalRouteHandler.addSearch) // ADMIN ROUTE COMMENTED TEMP
		// loggedIn check needed
		global.POST("/upload-file", GlobalRouteHandler.imageRouteHandler)
		//router.GET("/discord/authorize", GlobalRouteHandler.authorizeDiscord)
		//studio.PUT("/members/:userId", member.MemberRouteHandler.AddMemberToStudio)
	}
	integrations := router.Group("integrations")
	{
		// Only for studio integration => check for STUDIO_MANAGE_INTEGRATION
		// @todo Will add this later once discord is stable
		// TODO to be moved to studio integrations app
		integrations.GET("/discord/connect", GlobalRouteHandler.connectDiscord)
	}
	message := router.Group("message")
	{
		message.GET("/get", GlobalRouteHandler.getMessages)
		message.DELETE("/:messageID", GlobalRouteHandler.deleteMessage)
	}

}

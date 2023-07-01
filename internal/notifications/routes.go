package notifications

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a NotificationApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("notifications")
	{
		// loggedin
		App.Routes.GET("", App.RouteHandler.getNotifications)
		//App.Routes.POST("/mark-as-read", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.markAsRead)
		App.Routes.GET("/count", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.notificationCount)
		App.Routes.PATCH("/update", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.updateNotification)
		// self
		App.Routes.POST("/mark-as-seen", App.RouteHandler.markAsSeen)
		//App.Routes.GET("/count", App.RouteHandler.notificationCount)
		//App.Routes.PATCH("/update", App.RouteHandler.updateNotification)
	}
}

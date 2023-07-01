package studio_integration

import "github.com/gin-gonic/gin"

func (a StudioIntegrationApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("integrations")
	{
		App.Routes.GET("/settings", App.RouteHandler.GetSettings)
		App.Routes.PUT("/discord", App.RouteHandler.UpdateDiscordDmNotifications)
		App.Routes.PUT("/slack", App.RouteHandler.UpdateSlackDmNotifications)
		App.Routes.DELETE("", App.RouteHandler.DeleteIntegration)
		App.Routes.GET("/discord/update", App.RouteHandler.CheckUpdateIntegration)
	}
}

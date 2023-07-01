package slack2

import (
	"github.com/gin-gonic/gin"
	// "golang.org/x/oauth2/slack"
)

// WebhookImpl has all methods
// Update the Url for prod in slack apps
// 1. Interactivity & Shortcuts with /slack/shortcuts
// 2. Event Subscriptions with /slack/events
// Update bot subscription events app_mention,
// https://api.dev.bip.so/api/v1/slack/events

func (a slackApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("slack")
	{
		App.Routes.GET("/connect", App.RouteHandler.connectSlack)
		App.Routes.GET("/authorize", App.RouteHandler.authorize)
		App.Routes.POST("/connect_login", App.RouteHandler.connectSlackLogin)

		// Added under Interactivity & Shortcuts in slack apps https://api.slack.com/
		App.Routes.POST("/shortcuts", App.RouteHandler.SlackShortcutsHandler)
		App.Routes.POST("/events", App.RouteHandler.SlackEventsHandler)
		App.Routes.POST("/slash-commands", App.RouteHandler.SlashCommandsHandler)
	}
}

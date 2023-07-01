package discord

import (
	"github.com/gin-gonic/gin"
)

type DiscordImpl struct{}

func RegisterRoutes(router *gin.RouterGroup) {
	impl := &DiscordImpl{}
	// router.POST("/discord/connect_login", impl.connectDiscordLogin)
	// router.GET("/discord/authorize", impl.authorizeDiscord)
	// 	This end point is set on the discord side
	// https://discord.com/developers/applications/856212411908751411/information
	// Current it is set to : https://api-v2-eu.prod-deployment.bip.so/discord/interaction
	// Slash Commands
	router.POST("/discord/interaction", impl.interaction)
}

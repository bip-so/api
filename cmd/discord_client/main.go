package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/cmd/discord_client/connection"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
)

func appSetup() {
	// Init config from the env file.
	configs.InitConfig(".env", ".")
	// Start everything : Logger / DB / Redis : Need Err.
	core.InitCore(".env", ".")
}
func main() {
	fmt.Println("⇾ Init Discord Client.")
	appSetup()
	r := gin.Default()
	fmt.Println("⇾ Starting Bot Connection...")
	// Starting
	connection.StartDiscordBotConnection(configs.GetConfigString("DISCORD_TOKEN"), "bip")
	r.Run("localhost:9002")
}

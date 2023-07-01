package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/cmd/discord_worker/eventHandler"
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
	fmt.Println("⇾ Init Discord Worker.")
	appSetup()
	r := gin.Default()
	fmt.Println("⇾ Processing Events")
	eventHandler.Consumer()
	r.Run("localhost:9003")
}

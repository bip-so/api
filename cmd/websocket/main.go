package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/phonepost/bip-be-platform/cmd/api"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"net/http"
)

func appSetup() {
	//fileName := ".env"
	fileName := ".env"
	// Init config from the env file.
	configs.InitConfig(fileName, ".")
	// Start everything : Logger / DB / Redis : Need Err.
	core.InitCore(fileName, ".")
	api.InitAllApps()
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	fmt.Println("Connected")
	defer conn.Close()
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}
func main() {
	appSetup()
	r := gin.Default()
	//r.GET("/", func(c *gin.Context) {
	//	c.String(200, "Websocket !")
	//})
	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	r.Run("localhost:9009")
}

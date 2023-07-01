package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func RegisterBootstrapRoutes(router *gin.RouterGroup) {
	bootstrap := router.Group("bootstrap")
	{
		bootstrap.GET("/get", middlewares.TokenAuthorizationMiddleware(), bootstrapRouteHandler.Get)
		bootstrap.GET("/handle/:handle", bootstrapRouteHandler.Handle)
		bootstrap.GET("/user/:userid",bootstrapRouteHandler.Getuser)
	}
}

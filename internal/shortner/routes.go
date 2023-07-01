package shortner

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a ShortApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("short")
	{
		App.Routes.GET("/:shortID", App.RouteHandler.Get)
		App.Routes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Create)
	}
}

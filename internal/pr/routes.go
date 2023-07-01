package pr

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a PrApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("publish-requests")
	{
		App.Routes.GET("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.ListPublishRequests)

	}
}

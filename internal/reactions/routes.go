package reactions

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a ReactionApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("reactions")
	{
		// Add reaction -> Block -> CANVAS_BRANCH_ADD_REACTION
		// Add reaction -=>  Reel -> CANVAS_BRANCH_REACT_TO_REEL
		// Add reaction -> Block Comment -> CANVAS_BRANCH_ADD_REACTION
		// Add reaction -> Reel Comment -> CANVAS_BRANCH_REACT_TO_REEL
		App.Routes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Create)
		// 1.3 Get all Reactions on a Branch with block ID and UUID
		// self user reaction
		App.Routes.POST("/remove", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Remove)

	}
}

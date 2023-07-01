package mentions

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a MentionsApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("mentions")
	{
		// LoggedIn User
		/*
				CANVAS_BRANCH_EDIT -> block
				CANVAS_BRANCH_ADD_COMMENT -> block thread & block comment
			 	CANVAS_BRANCH_CREATE_REEL -> reel
				CANVAS_BRANCH_COMMENT_ON_REEL -> reel comment
		*/
		App.Routes.POST("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.AddMention)
	}
}

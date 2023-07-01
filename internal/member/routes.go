package member

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a MemberApp) RegisterRoutes(router *gin.RouterGroup) {
	App.Routes = *router.Group("member")
	{
		// self check
		App.Routes.POST(
			"/leave-studio/:studioId",
			middlewares.TokenAuthorizationMiddleware(),
			App.RouteHandler.LeaveStudio,
		)
		// CANVAS_BRANCH_VIEW
		App.Routes.GET("/canvas-branch/:canvasBranchID", App.RouteHandler.GetCanvasBranchMembers)
		// For now no specific check
		App.Routes.GET("/role/:roleID", App.RouteHandler.GetRoleMembers)
		// For now no specific check
		App.Routes.GET("/search", App.RouteHandler.SearchMembers)
		App.Routes.GET("/role/:roleID/search-members", App.RouteHandler.GetRoleMembersSearch)
	}
}

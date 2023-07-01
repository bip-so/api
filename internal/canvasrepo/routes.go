package canvasrepo

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a CanvasRepoApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("canvas-repo")
	{
		// This is a special api which creates a fresh Repo and inits a main branch on it,
		// You'll essentially get a CanvasRepo and

		// Move FE from Init to create
		// Immediate parent perms CANVAS_BRANCH_VIEW_METADATA
		App.Routes.POST("/init", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Init)
		// This is buggy needs to be Refactored
		App.Routes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Create)

		// No specific check
		App.Routes.POST("/get", App.RouteHandler.GetAllCanvas)
		// CANVAS_BRANCH_EDIT_NAME
		App.Routes.PATCH("/:canvasRepoID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Update)
		// CANVAS_BRANCH_VIEW
		App.Routes.POST("/create-language", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.CreateLanguage)
		// CANVAS_BRANCH_DELETE
		App.Routes.DELETE("/:canvasRepoID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.DeleteCanvasRepo)
		// STUDIO_CHANGE_CANVAS_COLLECTION_POSITION
		App.Routes.POST("/move", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.MoveCanvas)
		// get Canvas By key
		// CANVAS_BRANCH_VIEW_METADATA
		App.Routes.GET("/", App.RouteHandler.GetCanvas)
		// STUDIO_ADD_REMOVE_USER_TO_ROLE
		App.Routes.POST("/user/:userId", App.RouteHandler.GetMemberCanvas)
		App.Routes.POST("/role/:roleId", App.RouteHandler.GetRoleCanvas)
		App.Routes.GET("/search/user/:userId", App.RouteHandler.UserSearchCanvases)
		App.Routes.GET("/search/role/:roleId", App.RouteHandler.RoleSearchCanvases)

		App.Routes.GET("/next-prev/:canvasRepoId", App.RouteHandler.GetNextPrevCanvas)
		App.Routes.GET("/lang-next-prev/:canvasRepoId", App.RouteHandler.GetLangNextPrevCanvas)

		App.Routes.GET("/distinct-languages", App.RouteHandler.GetDistinctRepoLanguages)
	}
}

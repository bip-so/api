package role

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func RegisterRoleCrudRoutes(router *gin.RouterGroup) {
	roleCrud := router.Group("role").Use(middlewares.TokenAuthorizationMiddleware())
	{
		// STUDIO_CREATE_DELETE_ROLE, cannot delete or update name of bip-admin, member role
		roleCrud.POST("/create", middlewares.TokenAuthorizationMiddleware(), RoleCrudRouteHandler.createRole)    // Create Role
		roleCrud.POST("/edit", middlewares.TokenAuthorizationMiddleware(), RoleCrudRouteHandler.editRole)        // Edit Role
		roleCrud.DELETE("/:roleId", middlewares.TokenAuthorizationMiddleware(), RoleCrudRouteHandler.deleteRole) // Delete Role
		// Invalidate cache here too
		// STUDIO_ADD_REMOVE_USER_TO_ROLE
		roleCrud.POST("/membership", middlewares.TokenAuthorizationMiddleware(), RoleCrudRouteHandler.updateMembership) // Add or Remove Members
		// no specific check
		roleCrud.GET("/member/:memberId", RoleCrudRouteHandler.getMemberRoles)
	}
}

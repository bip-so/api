package studio

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a StudioApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("studio")
	{
		App.Routes.GET("/:studioId", App.RouteHandler.getStudioRouteHandler)
		// STUDIO_EDIT_STUDIO_PROFILE
		App.Routes.POST("/edit", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.editRouteHandler)
		App.Routes.GET("/toggle-membership", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.toggleStudioMembershipRouteHandler)
		// STUDIO_EDIT_STUDIO_PROFILE
		App.Routes.POST("/image", App.RouteHandler.imageRouteHandler)
		// no check
		App.Routes.GET("/popular", App.RouteHandler.popularStudiosRouteHandler)
		//studio.POST("/studiolist", App.RouteHandler.studioListRouteHandler)
		// STUDIO_DELETE
		App.Routes.DELETE("/:studioId", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.deleteRouteHandler)
		App.Routes.GET("/admins", App.RouteHandler.getStudioAdminMembers)
		App.Routes.GET("/members", App.RouteHandler.getStudioMembersRouteHandler)
		App.Routes.GET("/roles", App.RouteHandler.getStudioRolesRouteHandler)
		// STUDIO_CREATE_DELETE_ROLE
		App.Routes.POST("/ban", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.banUser)
		// self or STUDIO_MANAGE_PERMS
		App.Routes.POST("/:studioId/join", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.joinStudio)
		// STUDIO_MANAGE_PERMS
		App.Routes.POST("/join/bulk", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.joinStudioInBulk)
		// no check
		App.Routes.POST("/memberCount", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.memberCount)

		App.Routes.POST("/invite-via-email", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.InviteWithEmailFlow)

		// RTJ: get list
		App.Routes.GET("/:studioId/membership-request/list", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.requestToJoinStudioList)

		// RTJ: create request
		App.Routes.POST("/:studioId/membership-request/new", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.requestToJoinStudioCreate)
		// RTJ: reject
		App.Routes.POST("/:studioId/membership-request/:membershipRequestID/reject", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.requestToJoinStudioReject)
		// RTJ: accept
		App.Routes.POST("/:studioId/membership-request/:membershipRequestID/accept", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.acceptToJoinStudioReject)

		App.Routes.GET("/:studioId/stats", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.getStudioStats)

		// Also allows for the ?upgrade=true
		App.Routes.GET("/:studioId/customer-portal-session", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.NewCustomerPortalSession)
		App.Routes.GET("/:studioId/customer-payment-session", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.NewPaymentLinkSession)

		studioExternalIntegrationRoutes := App.Routes.Group("/external-integration")
		studioExternalIntegrationRoutes.GET("/ping", middlewares.VendorValidationMiddleware(), App.RouteExternalHandler.Ping)
		studioExternalIntegrationRoutes.GET("/", middlewares.VendorValidationMiddleware(), App.RouteExternalHandler.Get)
		studioExternalIntegrationRoutes.POST("/", middlewares.VendorValidationMiddleware(), App.RouteExternalHandler.Create)
		studioExternalIntegrationRoutes.GET("/user-points", middlewares.VendorValidationMiddleware(), App.RouteExternalHandler.GetUserPoints)
	}

}

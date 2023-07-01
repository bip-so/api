package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a AuthApp) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("auth")
	{
		// Anonymous Calls
		auth.POST("/signup", App.RouteHandler.signupLegacy)
		auth.POST("/login", App.RouteHandler.loginLegacy)
		// Forgot password: user has forgotten password, Need an email
		// This API will send an email if  exists with a simple token.
		auth.POST("/forgot-password", App.RouteHandler.forgotPasswordInit)
		auth.POST("/reset-password", App.RouteHandler.resetPassword)
		auth.GET("/ghost-login", App.RouteHandler.SpecialGhostLogin)
		// to do fix
		//auth.POST("/reset-password", App.RouteHandler.resetPassword)
		auth.POST("/change-password", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.changePassword)

		// Reset  password

		auth.POST("/existing-email", App.RouteHandler.existingEmail)
		auth.POST("/existing-username", App.RouteHandler.existingUsername)
		auth.POST("/otp", App.RouteHandler.generateUserOtp)
		//http://localhost:9001/api/v1/auth/verify-email/1810859c-8686-47e8-9a00-998c8a9f900d -> UUID is user's UUID
		auth.GET("/verify-email/:verificationKey", App.RouteHandler.VerifyUserEmail)

		// Auth required
		auth.POST("/logout", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.logout)
		// Todo: We need a new middle ware to check the refresh token insted of tghe access token
		auth.POST("/refresh-token", App.RouteHandler.refreshToken)
		auth.POST("/social-login", App.RouteHandler.SocialLogin)

	}
}

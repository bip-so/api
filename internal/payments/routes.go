package payments

import (
	"github.com/gin-gonic/gin"
)

type StripeImpl struct{}

func RegisterStripeRoutes(router *gin.RouterGroup) {
	stripeRouter := router.Group("stripe")
	{
		stripeRouter.POST("/create-checkout-session", App.Controller.CreateCheckoutSession)
		stripeRouter.POST("/success", App.Controller.Success)
		stripeRouter.POST("/cancel", App.Controller.Cancel)
	}
}

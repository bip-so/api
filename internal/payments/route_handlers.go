package payments

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/checkout/session"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

func (controller paymentController) CreateCheckoutSession(c *gin.Context) {
	// Get or create a customer in our pg db by mapping of userID and stripe customerId and
	// pass the customer id while creating the session.
	// Based on the flow we can create StripeCustomer and use that table for mapping.
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Basic Plan"),
					},
					UnitAmount: stripe.Int64(2000),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(fmt.Sprintf("%s", configs.GetAppInfoConfig().FrontendHost)),
		CancelURL:  stripe.String(fmt.Sprintf("%s", configs.GetAppInfoConfig().FrontendHost)),
	}

	s, err := session.New(params)
	if err != nil {
		fmt.Println("errror", err)
		return
	}
	c.Redirect(http.StatusSeeOther, s.URL)
	return
}

func (controller paymentController) Success(c *gin.Context) {
	fmt.Println(c.Request.Body)
}

func (controller paymentController) Cancel(c *gin.Context) {
	fmt.Println(c.Request.Body)
}

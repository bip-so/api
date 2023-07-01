package payments

import (
	"github.com/stripe/stripe-go/v73"
	billingportalsession "github.com/stripe/stripe-go/v73/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v73/checkout/session"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/stripe/stripe-go/v73/subscription"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strconv"
)

func (s *paymentService) GetStripeData() map[string]string {
	ENV := configs.GetConfigString("ENV")
	return utils.STRIPE_CONST_LOOKUPS[ENV]
}

func (s *paymentService) AddSub(customerID string, priceID string, memberCount int64) string {
	sd := s.GetStripeData()
	stripe.Key = sd["SECRET_KEY"]

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(memberCount),
			},
		},
	}
	sub, _ := subscription.New(params)
	return sub.ID

}

func (s *paymentService) DeleteCustomerOnStripe(studioID uint64) (string, error) {
	//fmt.Println(userInstance)
	sd := s.GetStripeData()
	stripe.Key = sd["SECRET_KEY"]
	studioInstance, _ := queries.App.StudioQueries.GetStudioInstanceByID(studioID)
	//params := &stripe.CustomerParams{
	//	Name:  stripe.String(studioInstance.DisplayName),
	//	Email: stripe.String(email),
	//}
	//params.AddMetadata("studio_id", strconv.FormatUint(studioInstance.ID, 10))
	//params.AddMetadata("handle", studioInstance.Handle)
	//params.AddMetadata("display_name", studioInstance.DisplayName)
	//params.AddMetadata("clientReferenceId", userInstance.ClientReferenceId)
	c, _ := customer.Del(studioInstance.StripeSubscriptionsID, nil)
	//sub := s.AddSub(c.ID, sd["LITE_PRICE"], memberCount)
	return c.ID, nil
}

func (s *paymentService) CreateNewCustomerOnStripe(studioID uint64) (string, error) {
	//fmt.Println(userInstance)
	sd := s.GetStripeData()
	stripe.Key = sd["SECRET_KEY"]
	var email string
	studioInstance, err := queries.App.StudioQueries.GetStudioInstanceByID(studioID)
	userInstance, _ := queries.App.UserQueries.GetUser(map[string]interface{}{"id": studioInstance.CreatedByID})

	if studioInstance.CreatedByUser.Email.Valid {
		email = studioInstance.CreatedByUser.Email.String
	}
	if err != nil {
		return "", err
	}
	params := &stripe.CustomerParams{
		Name:  stripe.String(studioInstance.DisplayName),
		Email: stripe.String(email),
	}
	params.AddMetadata("studio_id", strconv.FormatUint(studioInstance.ID, 10))
	params.AddMetadata("handle", studioInstance.Handle)
	params.AddMetadata("display_name", studioInstance.DisplayName)
	params.AddMetadata("clientReferenceId", userInstance.ClientReferenceId)

	c, _ := customer.New(params)
	memberCount := queries.App.StudioQueries.StudioMemberCount(studioInstance.ID)
	//sub := s.AddSub(c.ID, sd["LITE_PRICE"], memberCount)

	_ = queries.App.StudioQueries.UpdateStudioByID(studioID, map[string]interface{}{
		"stripe_customer_id": c.ID,
		"stripe_product_id":  sd["LITE_PRODUCT_ID"],
		"stripe_price_id":    sd["LITE_PRICE"],
		"stripe_price_unit":  memberCount,
		//"stripe_subscriptions_id": sub,
	})

	return c.ID, nil
}

func (s *paymentService) PortalLink(studioID uint64, url string, loggedInUser *models.User) map[string]interface{} {
	var redirectUrl *string
	sd := s.GetStripeData()
	stripe.Key = sd["SECRET_KEY"]
	if url != "" {
		redirectUrl = &url
	} else {
		redirectUrl = stripe.String(sd["REDIRECT_URL"])
	}
	studioInstance, _ := queries.App.StudioQueries.GetStudioInstanceByID(studioID)

	// Handling Edge Case (If the customer has no Stripe Details.
	if studioInstance.StripeCustomerID == "na" {
		App.Service.CreateNewCustomerOnStripe(studioInstance.ID)
		studioInstance, _ = queries.App.StudioQueries.GetStudioInstanceByID(studioID)
	}

	// We can check if customer has a plan or not if they don't have a plan we send them the payment link
	hasProSub := s.hasProSub(&studioInstance.StripeCustomerID)
	if !hasProSub {
		return s.CheckoutSession(studioID, url)
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(studioInstance.StripeCustomerID),
		ReturnURL: redirectUrl,
	}
	//params.AddMetadata()
	sess, _ := billingportalsession.New(params)
	stripeSessionURL := sess.URL

	//if upgrade {
	//	redirectUrl = sess.URL + "/subscriptions/" + studioInstance.StripeSubscriptionsID + "/preview/" + sd["PRO_PRICE"]
	//}
	// Is user on a lower plan ?
	//if studioInstance.StripePriceID == sd["LITE_PRICE"] {
	//	stripeSessionURL = sess.URL + "/subscriptions/" + studioInstance.StripeSubscriptionsID + "/preview/" + sd["PRO_PRICE"]
	//	return map[string]interface{}{
	//		"stripe_session_url": sess.URL,
	//		"url":                stripeSessionURL,
	//		"return_url":         &redirectUrl,
	//	}
	//}

	return map[string]interface{}{
		//"stripe_session_url": sess.URL, -> Need
		"url":        stripeSessionURL,
		"return_url": &redirectUrl,
	}
}

func (s *paymentService) hasProSub(stripeCustomerID *string) bool {
	sd := s.GetStripeData()
	stripe.Key = sd["SECRET_KEY"]

	params := &stripe.SubscriptionListParams{}
	params.Customer = stripeCustomerID
	params.Price = stripe.String(sd["PRO_PRICE"])
	params.Filters.AddFilter("limit", "", "3")
	i := subscription.List(params)
	hasProSub := false

	for i.Next() {
		//s1 := i.Subscription()
		hasProSub = true
		//fmt.Println("SubID", s1.ID)
		//fmt.Println("customerID", s1.Customer)
		////sub_1M9ivuSF3qmDjcKMMhYivpFH

	}
	return hasProSub
}

func (s *paymentService) CheckoutSession(studioID uint64, url string) map[string]interface{} {
	var redirectUrl *string
	sd := s.GetStripeData()
	stripe.Key = sd["SECRET_KEY"]

	if url != "" {
		redirectUrl = &url
	} else {
		redirectUrl = stripe.String(sd["REDIRECT_URL"])
	}
	studioInstance, _ := queries.App.StudioQueries.GetStudioInstanceByID(studioID)

	userInstance, _ := queries.App.UserQueries.GetUser(map[string]interface{}{"id": studioInstance.CreatedByID})

	// Handling Edge Case (If the customer has no Stripe Details.
	if studioInstance.StripeCustomerID == "na" {
		App.Service.CreateNewCustomerOnStripe(studioInstance.ID)
		studioInstance, _ = queries.App.StudioQueries.GetStudioInstanceByID(studioID)
	}

	var paymentSession *stripe.CheckoutSession
	var params2 *stripe.CheckoutSessionParams
	if userInstance.ClientReferenceId == "na" {
		params2 = &stripe.CheckoutSessionParams{
			Customer:            stripe.String(studioInstance.StripeCustomerID),
			SuccessURL:          redirectUrl,
			CancelURL:           redirectUrl,
			AllowPromotionCodes: stripe.Bool(true),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				&stripe.CheckoutSessionLineItemParams{
					Price:    stripe.String(sd["PRO_PRICE"]),
					Quantity: stripe.Int64(studioInstance.StripePriceUnit),
				},
			},
			Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		}
		paymentSession, _ = checkoutsession.New(params2)
	} else {
		params2 = &stripe.CheckoutSessionParams{
			SuccessURL: redirectUrl,
			CancelURL:  redirectUrl,
			Customer:   stripe.String(studioInstance.StripeCustomerID),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				&stripe.CheckoutSessionLineItemParams{
					Price:    stripe.String(sd["PRO_PRICE"]),
					Quantity: stripe.Int64(studioInstance.StripePriceUnit),
				},
			},
			Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			ClientReferenceID: stripe.String(userInstance.ClientReferenceId),
		}
		paymentSession, _ = checkoutsession.New(params2)
	}

	return map[string]interface{}{
		"url":        paymentSession.URL,
		"return_url": &redirectUrl,
	}
}

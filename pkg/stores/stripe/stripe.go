package stripe

import (
	"github.com/stripe/stripe-go/v73"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

func InitStripe() {
	stripe.Key = configs.GetConfigString("STRIPE_SECRET_KEY")
}

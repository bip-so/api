package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"net/http"
)

// bip-partner-key
// startit
// stage : 'xNLQq3qcifD59pSnHqxJIA'
// prod: 'vBLJxDE5nSmu9O6vbQ4sbA'
const VendorHeaderKeyName = "Bip-Partner-Key"

func VendorValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		allowedVendorKeys := []string{"xNLQq3qcifD59pSnHqxJIA", "vBLJxDE5nSmu9O6vbQ4sbA"}
		bipVendorHeaderKey := c.Request.Header.Get(VendorHeaderKeyName)
		if bipVendorHeaderKey == "" {
			c.String(http.StatusUnauthorized, "No Vendor Authorization header provided")
			c.Abort()
			return
		}
		if !utils.SliceContainsItem(allowedVendorKeys, bipVendorHeaderKey) {
			c.String(http.StatusUnauthorized, "Incorrect Vendor Keys Provided")
			c.Abort()
			return
		}
		if bipVendorHeaderKey == "xNLQq3qcifD59pSnHqxJIA" || bipVendorHeaderKey == "vBLJxDE5nSmu9O6vbQ4sbA" {
			c.Set("vendorName", "start_it")
		}
		c.Next()
	}
}

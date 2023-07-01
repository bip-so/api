package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// Constraint process OPTIONS request.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := c.Request.Header.Get("Origin"); origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,UniqueUserID,bip-client-id,bip-studio-id, bip-token-invalid")
			c.Writer.Header().Add("Access-Control-Expose-Headers", "bip-client-id,bip-studio-id, bip-token-invalid")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,PATCH,DELETE")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}
		if c.Request.Method == "OPTIONS" {
			response.RenderBlankResponse(c)
			return
		}
		c.Next()
	}
}

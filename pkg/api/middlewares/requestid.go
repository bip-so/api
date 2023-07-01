package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strconv"
)

const headerKey = "X-Bip-Request-ID"

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(headerKey, utils.NewNanoid())
		var studioIdHeader uint64
		// Temp:
		studio_id := c.Request.Header.Get("bip-studio-id")
		if studio_id == "" {
			studioIdHeader = 999999999999
		} else {
			studioId, _ := strconv.ParseUint(studio_id, 10, 64)
			studioIdHeader = studioId
			c.Set("currentStudio", studioIdHeader)
		}

		c.Next()
	}
}

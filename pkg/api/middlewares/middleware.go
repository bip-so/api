// placeholder: To be done

package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// https://go.dev/blog/constants
const ClientID string = "bba53dba-ae8b-4890-a249-717bbcf48c3b"
const ClientHeaderKeyName = "Bip-Client-Id"

func ServerClientHeaderCkecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Hello world")
		mode := os.Getenv("APP_MODE")
		if mode == "local" {
			c.Next()
		}
		bipClientheader := c.Request.Header[ClientHeaderKeyName]
		ipfromBipClientheader := bipClientheader[0]
		if ipfromBipClientheader != ClientID {

			c.JSON(http.StatusUnauthorized, "Bad Client, Abort")
			c.Abort()
			return
		} else {
			c.Next()
		}

	}
}

func SetStudioInContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		studio_id := c.Request.Header.Get("bip-studio-id")

		if studio_id == "" {
			c.JSON(http.StatusNotFound, "StudioId not present")
			c.Abort()
			return
		}

		// studioId, _ := strconv.ParseUint(studio_id, 10, 64)
		// studio, err := studio.GetStudioByID(studioId)
		// if err != nil {
		// 	c.JSON(http.StatusNotFound, err.Error())
		// 	c.Abort()
		// 	return
		// }
		// c.Set("currentStudio", studio)
		// c.Next()
	}
}

func StudioHeaderRequiredCheck(c *gin.Context) (uint64, bool) {
	ctxStudio, _ := c.Get("currentStudio")
	ctxStudioAsInt := ctxStudio.(uint64)
	if ctxStudioAsInt == 0 {
		//response.RenderCustomResponse(c, map[string]interface{}{
		//	"error": "Studio Header Missing",
		//})
		return 0, false
	}

	return ctxStudioAsInt, true
}

func SentryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if configs.IsLive() || configs.IsDev() {
			defer func() {
				err := recover()
				if err != nil {
					sentry.CurrentHub().Recover(err)
					sentry.Flush(time.Second * 5)
					fmt.Println("Panic Error recoved by Sentry::", err)
					response.RenderInternalServerErrorResponse(c, err)
				}
			}()
			c.Next()
		}
	}
}

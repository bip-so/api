package twitter

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

type TwitterImpl struct{}

func RegisterRoutes(router *gin.RouterGroup) {
	impl := &TwitterImpl{}
	router.GET("/twitter/metadata/:id", middlewares.TokenAuthorizationMiddleware(), impl.GetMetadata)
}

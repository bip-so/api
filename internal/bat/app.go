package bat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

type batRoutes struct {
	core.RouteHelper
}

type batController struct {
	logg logger.BipLogger
}
type batService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type batApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler batRoutes
	Controller   batController // exposes logg
	Service      batService    // exposes logg and cache
}

var App batApp

func InitApp() {
	App.Name = "bat"
	App.Service.cache = bipredis.RedisClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

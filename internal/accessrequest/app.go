package ar

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

type arRoutes struct {
	core.RouteHelper
}

type arController struct {
	logg logger.BipLogger
}
type arService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type arApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler arRoutes
	Controller   arController // exposes logg
	Service      arService    // exposes logg and cache
}

var App arApp

func InitApp() {
	App.Name = "ar"
	App.Service.cache = bipredis.RedisClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

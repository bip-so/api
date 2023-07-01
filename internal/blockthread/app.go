package blockthread

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gorm.io/gorm"
)

type blockThreadRoutes struct {
	core.RouteHelper
}

type blockThreadController struct {
	logg logger.BipLogger
}
type blockThreadService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type blockThreadRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type blockThreadApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler blockThreadRoutes
	Controller   blockThreadController // exposes logg
	Service      blockThreadService    // exposes logg and cache
	Repo         blockThreadRepo       // exposes db and cache
}

var App blockThreadApp

func InitApp() {
	App.Name = "blockThread"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

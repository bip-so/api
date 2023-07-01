package reel

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gorm.io/gorm"
)

type reelRoutes struct {
	core.RouteHelper
}

type reelController struct {
	logg logger.BipLogger
}
type reelService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	scope   *gorm.DB
}

type reelCachingService struct {
	cache *bipredis.Cache
}

type reelRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type ReelApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler reelRoutes
	Controller   reelController // exposes logg
	Service      reelService    // exposes logg and cache
	Repo         reelRepo       // exposes db and cache
	Caching      reelCachingService
}

var App ReelApp

func InitApp() {
	App.Name = "Reels"
	App.Repo.db = postgres.GetDB()
	App.Service.scope = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	App.Caching.cache = bipredis.NewCache()

	fmt.Println(App.Name + " started. ")
}

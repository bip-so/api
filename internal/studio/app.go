package studio

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

type studioRoutes struct{ core.RouteHelper }
type studioExternalRoutes struct{ core.RouteHelper }

type studioController struct {
	logg logger.BipLogger
}
type studioService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	db      *gorm.DB
}
type userAssociatedStudiosService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type studioTopicService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type studioRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type studioTopicRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type userAssociatedStudioRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type StudioApp struct {
	Name                 string
	Routes               gin.RouterGroup
	RouteHandler         studioRoutes
	RouteExternalHandler studioExternalRoutes

	Controller                   studioController   // exposes logg
	StudioService                studioService      // exposes logg and cache
	TopicService                 studioTopicService // exposes logg and cache
	UserAssociatedStudiosService userAssociatedStudiosService
	StudioRepo                   studioRepo // exposes db and cache
	UserAssociatedStudioRepo     userAssociatedStudioRepo
	TopicRepo                    studioTopicRepo
}

var App StudioApp

func InitApp() {
	App.Name = "Studio"

	App.StudioRepo.db = postgres.GetDB()
	App.StudioRepo.cache = bipredis.RedisClient()
	App.StudioRepo.kafka = kafka.GetKafkaClient()

	App.UserAssociatedStudioRepo.db = postgres.GetDB()
	App.UserAssociatedStudioRepo.cache = bipredis.RedisClient()
	App.UserAssociatedStudioRepo.kafka = kafka.GetKafkaClient()

	App.TopicRepo.db = postgres.GetDB()
	App.TopicRepo.cache = bipredis.RedisClient()
	App.TopicRepo.kafka = kafka.GetKafkaClient()

	App.StudioService.cache = bipredis.RedisClient()
	App.StudioService.kafka = kafka.GetKafkaClient()
	App.StudioService.db = postgres.GetDB()
	App.TopicService.cache = bipredis.RedisClient()
	App.TopicService.kafka = kafka.GetKafkaClient()

	App.UserAssociatedStudiosService.cache = bipredis.RedisClient()
	App.UserAssociatedStudiosService.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

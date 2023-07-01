package studio_integration

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

type studioIntegrationsRoutes struct {
	core.RouteHelper
}

type studioIntegrationsController struct {
	logg logger.BipLogger
}
type studioIntegrationsService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	scope   *gorm.DB
}

type studioIntegrationsRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type StudioIntegrationApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler studioIntegrationsRoutes
	Controller   studioIntegrationsController // exposes logg
	Service      studioIntegrationsService    // exposes logg and cache
	Repo         studioIntegrationsRepo       // exposes db and cache
}

var App StudioIntegrationApp

func InitApp() {
	App.Name = "StudioIntegration"
	App.Repo.db = postgres.GetDB()
	App.Service.scope = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

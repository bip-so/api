package notion

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

type notionRoutes struct {
	core.RouteHelper
}

type notionController struct {
	logg logger.BipLogger
}
type notionService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	scope   *gorm.DB
}

type notionRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type NotionApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler notionRoutes
	Controller   notionController // exposes logg
	Service      notionService    // exposes logg and cache
	Repo         notionRepo       // exposes db and cache
}

var App NotionApp

func InitApp() {
	App.Name = "Notion"
	App.Repo.db = postgres.GetDB()
	App.Service.scope = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

package slack2

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

type slackRoutes struct {
	core.RouteHelper
}

type slackController struct {
	logg logger.BipLogger
}
type slackService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type slackRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type slackApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler slackRoutes
	Controller   slackController // exposes logg
	Service      slackService    // exposes logg and cache
	Repo         slackRepo       // exposes db and cache
}

var App slackApp

func InitApp() {
	App.Name = "Slack App"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

package mentions

// This package routes are on CanvasBranch
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

type mentionsRoutes struct {
	core.RouteHelper
}

type mentionsController struct {
	logg logger.BipLogger
}
type mentionsService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type mentionsRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type MentionsApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler mentionsRoutes
	Controller   mentionsController // exposes logg
	Service      mentionsService    // exposes logg and cache
	Repo         mentionsRepo       // exposes db and cache
}

var App MentionsApp

func InitApp() {
	App.Name = "mrs"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

package blockThreadCommentcomment

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

type blockThreadCommentRoutes struct {
	core.RouteHelper
}

type blockThreadCommentController struct {
	logg logger.BipLogger
}
type blockThreadCommentService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type blockThreadCommentRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type blockThreadCommentApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler blockThreadCommentRoutes
	Controller   blockThreadCommentController // exposes logg
	Service      blockThreadCommentService    // exposes logg and cache
	Repo         blockThreadCommentRepo       // exposes db and cache
}

var App blockThreadCommentApp

func InitApp() {
	App.Name = "blockThreadComment"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

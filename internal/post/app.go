package post

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

type postRoutes struct {
	core.RouteHelper
}

type postController struct {
	logg logger.BipLogger
}
type postService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	db      *gorm.DB
}
type postRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type postApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler postRoutes
	Controller   postController // exposes logg
	Service      postService    // exposes logg and cache
	Repo         postRepo       // exposes db and cache
}

var App postApp

func InitApp() {
	App.Name = "Posts"
	App.Repo.db = postgres.GetDB()
	App.Service.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

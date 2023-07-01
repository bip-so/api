package canvasrepo

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

type canvasRepoRoutes struct{ core.RouteHelper }

type canvasRepoController struct {
	logg    logger.BipLogger
	Manager core.QuerySet
}
type canvasRepoService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type canvasRepoRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	Manager core.QuerySet
	kafka   *kafka.BipKafka
}
type repoUserCachingService struct {
	cache       *bipredis.Cache
	redisClient *redis.Client
}
type repoAutoPublichCachingService struct {
	cache       *bipredis.Cache
	redisClient *redis.Client
}

type CanvasRepoApp struct {
	Name             string
	Routes           gin.RouterGroup
	RouteHandler     canvasRepoRoutes
	Controller       canvasRepoController // exposes logg
	Service          canvasRepoService    // exposes logg and cache
	Repo             canvasRepoRepo       // exposes db and cache
	UserRepoHistory  repoUserCachingService
	PlansAutoPublish repoAutoPublichCachingService
}

//
//type Self interface {
//	CurrentUserID(c *gin.Context) uint64
//}

var App CanvasRepoApp

func InitApp() {
	App.Name = "Canvas Repo"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	App.UserRepoHistory.redisClient = bipredis.RedisClient()
	App.PlansAutoPublish.redisClient = bipredis.RedisClient()
	App.PlansAutoPublish.cache = bipredis.NewCache()
	//App.UserRepoHistory.cache = bipredis.RedisClient()

	fmt.Println(App.Name + " started. ")
}

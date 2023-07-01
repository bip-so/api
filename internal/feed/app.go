package feed

import (
	"fmt"
	"github.com/GetStream/stream-go2/v7"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	bipStream "gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gorm.io/gorm"
)

type feedRoutes struct{ core.RouteHelper }

type feedController struct {
	logg logger.BipLogger
}
type feedService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	Stream  *stream.Client
}
type feedRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	Stream  *stream.Client
}

type feedApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler feedRoutes
	Controller   feedController // exposes logg
	Service      feedService    // exposes logg and cache
	Repo         feedRepo       // exposes db and cache
}

var App feedApp

func InitApp() {
	App.Name = "feed"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	App.Repo.Stream = bipStream.StreamClient()
	App.Service.Stream = bipStream.StreamClient()
	fmt.Println(App.Name + " started. ")
}

package payments

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

type paymentRoutes struct {
	core.RouteHelper
}

type paymentController struct {
	logg logger.BipLogger
}
type paymentService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type paymentRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type paymentApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler paymentRoutes
	Controller   paymentController // exposes logg
	Service      paymentService    // exposes logg and cache
	Repo         paymentRepo       // exposes db and cache
}

var App paymentApp

func InitApp() {
	App.Name = "Stripe App"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

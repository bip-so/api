package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

type authRoutes struct {
	core.RouteHelper
}

type authController struct {
	logg logger.BipLogger
}
type authService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type authCreateStudioService struct {
	logg logger.BipLogger
}

type authRepo struct {
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type AuthApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler authRoutes
	Controller   authController // exposes logg
	Service      authService    // exposes logg and cache
	Repo         authRepo       // exposes db and cache
	PostSignup   authCreateStudioService
}

var App AuthApp

func InitApp() {
	App.Name = "Auth"
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

package user

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

type userRoutes struct{ core.RouteHelper }

type userController struct {
	logg logger.BipLogger
}
type userService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type userRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type userFollowRepo struct{}

type UserApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler userRoutes
	Controller   userController // exposes logg
	Service      userService    // exposes logg and cache
	Repo         userRepo       // exposes db and cache
	FollowRepo   userFollowRepo
}

var App UserApp

func InitApp() {
	App.Name = "User"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

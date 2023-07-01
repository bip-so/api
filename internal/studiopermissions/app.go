package studiopermissions

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

type studioPermissionsRoutes struct{ core.RouteHelper }

type studioPermissionController struct {
	logg logger.BipLogger
}
type studioPermissionService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type studioPermissionRepo struct {
	db     *gorm.DB
	cache  *redis.Client
	kafka  *kafka.BipKafka
	Manger core.QuerySet
}

type StudioPermissionApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler studioPermissionsRoutes
	Controller   studioPermissionController // exposes logg
	Service      studioPermissionService    // exposes logg and cache
	Repo         studioPermissionRepo       // exposes db and cache
}

var App StudioPermissionApp

func InitApp() {
	App.Name = "StudioPermission"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

package collectionpermissions

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

type collectionPermissionsRoutes struct{ core.RouteHelper }

type collectionPermissionController struct {
	logg logger.BipLogger
}
type collectionPermissionService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type collectionPermissionRepo struct {
	db     *gorm.DB
	cache  *redis.Client
	kafka  *kafka.BipKafka
	Manger core.QuerySet
}

type CollectionPermissionApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler collectionPermissionsRoutes
	Controller   collectionPermissionController // exposes logg
	Service      collectionPermissionService    // exposes logg and cache
	Repo         collectionPermissionRepo       // exposes db and cache
}

var App CollectionPermissionApp

func InitApp() {
	App.Name = "CollectionPermission"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

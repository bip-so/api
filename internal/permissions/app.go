package permissions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gorm.io/gorm"
)

type permissionRoutes struct{ core.RouteHelper }

type permissionController struct {
	logg logger.BipLogger
}
type permissionService struct {
	logg    logger.BipLogger
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type permissionRepo struct {
	db     *gorm.DB
	cache  *bipredis.Cache
	kafka  *kafka.BipKafka
	Manger core.QuerySet
}

type permissionApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler permissionRoutes
	Controller   permissionController // exposes logg
	Service      permissionService    // exposes logg and cache
	Repo         permissionRepo       // exposes db and cache
}

var App permissionApp

func InitApp() {
	App.Name = "permission"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.NewCache()
	App.Service.cache = bipredis.NewCache()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

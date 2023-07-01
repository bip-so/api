package xpcontribs

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

type xpcontribRoutes struct {
	core.RouteHelper
}

type xpcontribController struct {
	logg logger.BipLogger
}
type xpcontribService struct {
	logg    logger.BipLogger
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type xpcontribRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type xpcontribApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler xpcontribRoutes
	Controller   xpcontribController // exposes logg
	Service      xpcontribService    // exposes logg and cache
	Repo         xpcontribRepo       // exposes db and cache
}

var App xpcontribApp

const MainStudioLogNameSpace = "studio-blocks-log:"
const MainStudioAuditLogNameSpace = "studio-editors-log:"
const MainStudioPointsNameSpace = "studio-user-points:"

func InitApp() {
	App.Name = "XP Contrib"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	//App.Service.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.NewCache()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

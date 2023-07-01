package canvasbranch

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

type canvasBranchRoutes struct{ core.RouteHelper }

type canvasBranchController struct {
	logg logger.BipLogger
}
type canvasBranchService struct {
	logg    logger.BipLogger
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	db      *gorm.DB
}
type gitService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type canvasBranchRepo struct {
	db      *gorm.DB
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	logg    logger.BipLogger
}

type CanvasBranchApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler canvasBranchRoutes
	Controller   canvasBranchController // exposes logg
	Service      canvasBranchService    // exposes logg and cache
	Repo         canvasBranchRepo       // exposes db and cache
	Git          gitService
}

var App CanvasBranchApp

func InitApp() {
	App.Name = "Canvas Branch"
	App.Repo.db = postgres.GetDB()
	App.Service.db = postgres.GetDB()
	App.Repo.cache = bipredis.NewCache()
	App.Service.cache = bipredis.NewCache()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

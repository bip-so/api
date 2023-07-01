package collection

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/pkg/core"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gorm.io/gorm"
)

type collectionRoutes struct{ core.RouteHelper }

type collectionController struct {
	logg logger.BipLogger
}
type collectionService struct {
	logg    logger.BipLogger
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type collectionRepo struct {
	db     *gorm.DB
	cache  *bipredis.Cache
	kafka  *kafka.BipKafka
	Manger core.QuerySet
}

type CollectionApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler collectionRoutes
	Controller   collectionController // exposes logg
	Service      collectionService    // exposes logg and cache
	Repo         collectionRepo       // exposes db and cache
}

var App CollectionApp

func InitApp() {
	App.Name = "Collection"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.NewCache()
	App.Service.cache = bipredis.NewCache()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

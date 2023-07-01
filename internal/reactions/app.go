package reactions

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

type reactionRoutes struct {
	core.RouteHelper
}

type reactionController struct {
	logg logger.BipLogger
}
type reactionService struct {
	logg    logger.BipLogger
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type reactionRepo struct {
	db      *gorm.DB
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type ReactionApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler reactionRoutes
	Controller   reactionController // exposes logg
	Service      reactionService    // exposes logg and cache
	Repo         reactionRepo       // exposes db and cache
}

var App ReactionApp

func InitApp() {
	App.Name = "Reaction"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.NewCache()
	App.Service.cache = bipredis.NewCache()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

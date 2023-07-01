package notifications

import (
	"fmt"
	"github.com/GetStream/stream-go2/v7"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	bipStream "gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gorm.io/gorm"
)

type notificationRoutes struct{ core.RouteHelper }

type notificationController struct {
}

type notificationService struct {
	logg    logger.BipLogger
	cache   *bipredis.Cache
	kafka   *kafka.BipKafka
	Manager core.QuerySet
	Stream  *stream.Client
}
type notificationRepo struct {
	db     *gorm.DB
	cache  *bipredis.Cache
	kafka  *kafka.BipKafka
	Manger core.QuerySet
	Stream *stream.Client
}

type NotificationApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler notificationRoutes
	Service      notificationService // exposes logg and cache
	Repo         notificationRepo    // exposes db and cache
	Controller   notificationController
}

var App NotificationApp

func InitApp() {
	App.Name = "Notification"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.NewCache()
	App.Service.cache = bipredis.NewCache()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	App.Repo.Stream = bipStream.StreamClient()
	App.Service.Stream = bipStream.StreamClient()
	fmt.Println(App.Name + " started. ")
}

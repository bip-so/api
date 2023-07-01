package parser2

// This is a parser for the blocks on a branch
// We Get a Branch
// Get All Blocks

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gorm.io/gorm"
)

type Parser2 struct {
	BranchID uint64
	Blocks   *[]models.Block
}
type Command interface {
	PlainText()
	RichText()
}

type Utils struct{}

type parser2Routes struct {
	core.RouteHelper
}

type parser2Controller struct {
	logg logger.BipLogger
}
type parser2Service struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type parser2Repo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type Parser2App struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler parser2Routes
	Controller   parser2Controller // exposes logg
	Service      parser2Service    // exposes logg and cache
	Repo         parser2Repo       // exposes db and cache
}

var App Parser2App

func InitApp() {
	App.Name = "Parser2"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

package pr

/*

PERMISSIONS CAN BE ON EITHER COLLECTION OR CANVAS REPO -> DEFAULT BRANCH ONLY.,
You can only Ask for PR on a Default Branch or Something which is "default" inside Repo.
---

--
Perms to Check



NR's Notes : via Discord
So this have two phases
Phase 1:
	We check the CanvasBrach Type is it "default" based on CanvasRepo
	We branch is Default we will Check Perms On Branch
	Else we will only Check the Parent?
Phase 2;
	We will also check the permission on the "Collection" this canvas belongs to

---
 Till PR a Branch does not gett any Git functions.
---
So when user creates a new canvas, they need to publish it to the studio before they invite others etc.
Also, there is no git workflow in unpublished canvases today and the creator just writes directly in main branch.

*/

// This package routes are on CanvasBranch
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

type prRoutes struct {
	core.RouteHelper
}

type prController struct {
	logg logger.BipLogger
}
type prService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}
type prRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	kafka   *kafka.BipKafka
	Manager core.QuerySet
}

type PrApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler prRoutes
	Controller   prController // exposes logg
	Service      prService    // exposes logg and cache
	Repo         prRepo       // exposes db and cache
}

var App PrApp

func InitApp() {
	App.Name = "prs"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
	App.Repo.kafka = kafka.GetKafkaClient()
	App.Service.kafka = kafka.GetKafkaClient()
	fmt.Println(App.Name + " started. ")
}

package canvasbranchpermissions

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gorm.io/gorm"
)

type canvasBranchPermissionsRoutes struct{ core.RouteHelper }
type canvasBranchPermissionsController struct {
	logg logger.BipLogger
}
type canvasBranchPermissionsService struct {
	logg    logger.BipLogger
	cache   *redis.Client
	Manager core.QuerySet
}
type canvasBranchPermissionRepo struct {
	db      *gorm.DB
	cache   *redis.Client
	Manager core.QuerySet
}

type CanvasBranchPermissionsApp struct {
	Name         string
	Routes       gin.RouterGroup
	RouteHandler canvasBranchPermissionsRoutes
	Controller   canvasBranchPermissionsController // exposes logg
	Service      canvasBranchPermissionsService    // exposes logg
	Repo         canvasBranchPermissionRepo        // exposes db and cache
}

var App CanvasBranchPermissionsApp

func InitApp() {
	App.Name = "CanvasBranch Permissions"
	App.Repo.db = postgres.GetDB()
	App.Repo.cache = bipredis.RedisClient()
	App.Service.cache = bipredis.RedisClient()
}

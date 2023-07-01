package main

import (
	"gitlab.com/phonepost/bip-be-platform/cmd/api"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

func main() {
	configs.InitConfig(".env", ".")
	postgres.InitDB()
	integrations.InitDiscordGo()
	redis.InitRedis()
	api.InitAllApps()
	core.InitCore(".env", ".")
	//canvasrepo.App.Service.SendCollectionTreeToDiscord(257)
	var user *models.User
	postgres.GetDB().Where("id = ?", 82).First(&user)
	//global.CreateDefaultDocsInStudio(587, user)
}

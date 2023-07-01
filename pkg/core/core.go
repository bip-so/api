package core

import (
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/search"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/stripe"
)

func InitCore(configName, configPath string) {
	configs.InitConfig(configName, configPath)
	logger.InitLogger()
	postgres.InitDB()
	redis.InitRedis()
	kafka.InitKafka()
	search.InitAlgolia()
	integrations.InitDiscordGo()
	stream.InitStreamClient()
	apiClient.InitApiClient()
	stripe.InitStripe()
	// Pausing this
	//sentry.InitSentry()
}

func CoreHealth() error {
	db, err := postgres.GetDB().DB()
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

const BipLogoType = `
  ___ ___ ___ 
 | _ )_ _| _ \
 | _ \| ||  _/
 |___/___|_|  
`

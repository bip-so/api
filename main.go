package main

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/cmd/api"
	"gitlab.com/phonepost/bip-be-platform/internal/kafkatopics"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"log"

	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
)

func appSetup() {
	//envFilename := api.GetEnvFileName()
	envFilename := ".env"
	fmt.Print(envFilename)
	// Init config from the env file.
	//configs.InitConfig(".env", ".")
	configs.InitConfig(envFilename, ".")
	// Start everything : Logger / DB / Redis : Need Err.
	//core.InitCore(".env", ".")
	core.InitCore(envFilename, ".")
	api.InitAllApps()
	log.Println("----------------------------------\n\n")
}

// @title Bip Backend Platform
// @description  Bip Backend Platform server.
// @schemes http https
// @termsOfService https://bip.so/terms-of-service/

// @contact.name API Support
// @contact.url https://bip.so
// @contact.email santhosh@bip.so

// @license.name Apache 2.0
// @licence.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apiKey bearerAuth
// @in header
// @name Authorization
func main() {

	fmt.Print(core.BipLogoType)
	// config and dependencies
	appSetup()
	// DEBUG PURPOSES
	//fmt.Println("List of all env variables: ")
	//fmt.Println(os.Environ())
	// router
	if configs.GetConfigString(configs.KAFKA_CONSUMER_ENABLED) != "true" {
		api.StartHttpServer()
	} else {
		kafka.InitKafkaStartConsumer(kafkatopics.KafkaHandleTopics)
	}
}

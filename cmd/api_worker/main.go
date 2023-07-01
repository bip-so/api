package main

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/cmd/api"
	"gitlab.com/phonepost/bip-be-platform/cmd/discord_worker/eventHandler"
	"gitlab.com/phonepost/bip-be-platform/internal/tasks"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"log"

	"github.com/hibiken/asynq"
)

func appSetup() {
	//fileName := ".env"
	fileName := ".env"
	// Init config from the env file.
	configs.InitConfig(fileName, ".")
	// Start everything : Logger / DB / Redis : Need Err.
	core.InitCore(fileName, ".")
	eventHandler.InitDiscordGo()
	api.InitAllApps()
}

func main() {
	appSetup()
	redisAddr := fmt.Sprintf("%s:%s", configs.GetRedisConfig().Host, configs.GetRedisConfig().Port)
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: configs.GetRedisConfig().Password},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc("tasks:", tasks.DefaultTaskHandler)
	// ...register other handlers...
	go apiClient.AsyncScheduler()
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

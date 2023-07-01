package apiClient

import (
	"fmt"
	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"log"
)

const (
	DEFAULT     = "default"
	CRITICAL    = "critical"
	LOW         = "low"
	CommonRetry = 3
)

var client *asynq.Client

func InitApiClient() {
	redisAddr := fmt.Sprintf("%s:%s", configs.GetRedisConfig().Host, configs.GetRedisConfig().Port)
	client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: configs.GetRedisConfig().Password})
}

func AddToQueue(taskType string, payload []byte, queueName string, retry int) {
	task := asynq.NewTask(taskType, payload)
	if queueName == "" {
		queueName = "default"
	}
	info, err := client.Enqueue(task, asynq.Queue(queueName), asynq.MaxRetry(retry))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func GetApiClient() *asynq.Client {
	return client
}

package kafkatopics

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func HandleNotification(msg *kafka.Message) {
	var notification notifications.PostNotification
	err := json.Unmarshal(msg.Value, &notification)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}
	fmt.Println("REceived notification", msg.Key, notification)
	notifications.App.Service.CreateNotification(string(msg.Key), &notification)
}

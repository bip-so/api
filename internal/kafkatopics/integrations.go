package kafkatopics

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func HandleIntegrationsEvent(msg *kafka.Message) {
	if string(msg.Key) == models.REEL {
		var reel models.Reel
		err := json.Unmarshal(msg.Value, &reel)
		if err != nil {
			logger.Error(err.Error())
			KafkaConsumerError(msg, err)
			return
		}
		notifications.App.Service.SendToIntegration(&reel)
	} else if string(msg.Key) == models.POST {
		var post models.Post
		err := json.Unmarshal(msg.Value, &post)
		if err != nil {
			logger.Error(err.Error())
			KafkaConsumerError(msg, err)
			return
		}
		notifications.App.Service.SendPostToIntegration(&post)
	}
}

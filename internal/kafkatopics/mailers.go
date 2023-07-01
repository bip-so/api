package kafkatopics

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/mailers"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func HandleSendEmail(msg *kafka.Message) {
	var mail mailers.SendEmail
	err := json.Unmarshal(msg.Value, &mail)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}
	fmt.Println("Received email", msg.Key, mail)
	mailers.App.Service.ReceiveEmailEvent(string(msg.Key), mail)
}

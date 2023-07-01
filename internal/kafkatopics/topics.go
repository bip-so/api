package kafkatopics

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

/*
	Handling the kafka topics here.
	Here we are using switch cases to trigger different methods for different topics.
*/
func KafkaHandleTopics(msg *kafka.Message) {
	switch msg.Topic {
	case configs.KAFKA_TOPICS_NEW_USER:
		NewUserConsumer(msg)
		AddUserToFeedStream(msg)
	case configs.KAFKA_TOPICS_UPDATE_USER:
		UpdateUserConsumer(msg)
	case configs.KAFKA_TOPICS_NEW_STUDIO:
		NewStudioConsumer(msg)
	case configs.KAFKA_TOPICS_UPDATE_STUDIO:
		UpdateStudioConsumer(msg)
	case configs.KAFKA_TOPICS_DELETED_STUDIO:
		DeletedStudioConsumer(msg)
	case configs.KAFKA_TOPICS_NEW_CANVAS:
		UpdateCollectionCanvasCount(msg, true)
	case configs.KAFKA_TOPICS_DELETED_CANVAS:
		UpdateCollectionCanvasCount(msg, false)
	case configs.KAKFA_TOPICS_CALCULATE_PERMISSIONS:
		CalculatePermissions(msg)
	case configs.KAFKA_TOPICS_NOTIFICATIONS:
		HandleNotification(msg)
	case configs.KAFKA_TOPICS_EMAILS:
		HandleSendEmail(msg)
	case configs.KAFKA_INTEGRATION_EVENTS:
		HandleIntegrationsEvent(msg)
	default:
		logger.Debug(fmt.Sprintf("message at topic:%v partition:%v offset:%v	%s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value)))
	}
}

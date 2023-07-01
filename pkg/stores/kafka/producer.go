package kafka

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

/*
Publish a new event to the specific topic in kafka.

For publishing events to new topic:
	- We need to create the topic in the kafka using kafdrop here http://15.206.67.39:9000/
	- After creating the topic we can use it to publish the events from here.
	- To consume the events we need to add topic name to the KAFKA_CONSUMER_GROUP_TOPICS in kafka/const.go file
*/
func (k *BipKafka) Publish(topic string, key string, value []byte) {
	err := KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: value,
		},
	)

	if err != nil {
		logger.Error(err.Error())
	}
}

// KafkaWriter This is kafka core writer method. Instead of using the above publish method, we need to send an event to kafka we can use this method.
func (k *BipKafka) KafkaWriter() *kafka.Writer {
	return KafkaWriter
}

package kafka

import (
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

type BipKafka struct{}

var KafkaWriter *kafka.Writer

/*
	Initialize the kafka server. We start the producer and consuming server here.
	With Producer we can send events to the kafka.
	Consumer listens to all the events of subscribed Groups.
*/
func InitKafka() {
	kafkaConfig := configs.GetKafkaConfig()
	brokers := strings.Split(kafkaConfig.Hosts, ",")
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	KafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Balancer: &kafka.Hash{},
		Dialer:   dialer,
	})

	// started consuming the messages from kafka
	// go KafkaStartConsumer()
}

func GetKafkaClient() *BipKafka {
	return &BipKafka{}
}

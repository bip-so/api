package kafka

import (
	"context"
	"fmt"
	"strings"

	kafka "github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func getKafkaReader(kafkaURL string, groupTopics []string, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: groupTopics,
	})
}

/*
	Consumer get the data based on KAFKA_CONSUMER_GROUP_TOPICS, KAFKA_CONSUMER_GROUP_ID

	KAFKA_CONSUMER_GROUP_TOPICS: On every new topic we create on kafka we need to add it in the const file.
	KAFKA_CONSUMER_GROUP_ID:
		- Group Id is used to distribute the consumer events for the consumers which are present inside that group.
		For eg. If we spin up two servers of bip-be-platform server then there will be two consumers in same group and workload is
		transferred between two servers by round robin method.
*/
func InitKafkaStartConsumer(kafkaHandleTopics func(*kafka.Message)) {
	kafkaConfig := configs.GetKafkaConfig()
	kafkaUrl := fmt.Sprintf("%s", kafkaConfig.Hosts)
	reader := getKafkaReader(kafkaUrl, configs.KAFKA_CONSUMER_GROUP_TOPICS, configs.KAFKA_CONSUMER_GROUP_ID)

	defer reader.Close()

	logger.Debug("started consuming the kafka events ... !!")
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		kafkaHandleTopics(&msg)
	}
}

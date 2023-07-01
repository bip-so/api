package kafkatopics

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

const kafkaErrorRedisKey = "kafkaErrors:"

type kafkaRedisError struct {
	Error string
	Value []byte
}

/*
	Kafka consumer error stores the consumer errors in the redis cache.
*/
func KafkaConsumerError(msg *kafka.Message, err error) {
	redisKey := fmt.Sprintf("%s%s:%s", kafkaErrorRedisKey, msg.Topic, msg.Key)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	redisValue := kafkaRedisError{
		Error: err.Error(),
	}
	redisValueJson, _ := json.Marshal(redisValue)
	rc.Set(rctx, redisKey, redisValueJson, 0)
}

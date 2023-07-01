package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"os"
)

func setDataInRedis(event interface{}, eventID string, eventType string) {
	//cacheInstance := cache.NewCache()
	ctx := context.Background()
	keyData := []string{eventID}
	//cacheKey := cache.GenerateCacheKey("discordevents", keyData)
	cacheKey := RedisDiscordNamespace + keyData[0]
	render := EventGetSerializerData(event, eventType)
	dataStr, _ := json.Marshal(render)
	fmt.Println("Saving : ", cacheKey)
	//err := cacheInstance.Set(ctx, cacheKey, dataStr, &cache.Options{Expiration: time.Duration(7*24) * time.Hour}).Err()
	err := redis.RedisClient().Set(ctx, cacheKey, dataStr, 0).Err()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

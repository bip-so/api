package eventHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	cache "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

// consumes events from redis
//type eventMap map[string]int
//
//func (em eventMap) MarshalBinary() ([]byte, error) {
//	return json.Marshal(em)
//}

func Consumer() {
	cacheInstance := cache.NewCache()
	ctx := context.Background()
	InitDiscordGo()
	//eventMapValue := cacheInstance.Get(ctx, "discordevent:eventmap")
	//var eventMap map[string]int
	//if eventMapValue == nil {
	//	eventMap = map[string]int{}
	//	cacheInstance.Set(ctx, "discordevent:eventmap", eventMap, nil)
	//} else {
	//	eventMap = eventMapValue.(map[string]int)
	//}

	for {
		iter := cacheInstance.GetAllMatchingKeys(ctx, RedisDiscordNamespaceAll)
		for iter.Next(ctx) {
			value := cacheInstance.Get(ctx, iter.Val())
			keyStr := iter.Val()
			valStr := value.(string)
			fullKeySplit := strings.Split(keyStr, ":")
			keyName := fullKeySplit[1] // getting 2nd part
			fmt.Println(valStr)
			fmt.Println("Processing: ", keyName)
			//if eventMap[iter.Val()] >= 3 {
			//	cacheInstance.Delete(ctx, iter.Val())
			//delete(eventMap, iter.Val())
			//	cacheInstance.Set(ctx, "discordevent:eventmap", eventMap, nil)
			//	continue
			//}

			data := EventSerializer{}
			err := json.Unmarshal([]byte(valStr), &data)
			if err != nil {
				//eventMap[iter.Val()]++
				//cacheInstance.Set(ctx, "discordevent:eventmap", eventMap, nil)
				//continue
				// Do we need to keep a tab on Failed Tasks
				cacheInstance.Set(ctx, RedisDiscordNamespaceFailed+keyName, value, &cache.Options{Expiration: time.Hour * 168})
			}

			// process this value
			err = Process(data, valStr)
			if err != nil {
				//eventMap[iter.Val()]++
				//cacheInstance.Set(ctx, "discordevent:eventmap", eventMap, nil)
				//continue
				cacheInstance.Set(ctx, RedisDiscordNamespaceFailed+keyName, value, &cache.Options{Expiration: time.Hour * 24})
			}

			// success
			//cacheInstance.Set(ctx, RedisDiscordNamespaceProcessed+keyName, value, nil)
			cacheInstance.Delete(ctx, iter.Val())
			//delete(eventMap, iter.Val())
			//cacheInstance.Set(ctx, "discordevent:eventmap", eventMap, nil)
		}

	}

}

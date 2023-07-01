package post

import (
	"context"
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"log"
	"strconv"
)

// We need to store all the last x post for a studio

const POST_NS = "cached-studio-posts:"

func SetPostDataRedis(studioID uint64, data shared.PaginationData) {
	studioIDStr := strconv.Itoa(int(studioID))
	key := POST_NS + studioIDStr
	// Delete
	redis.RedisClient().Del(context.Background(), key)
	marshal, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
	}
	redis.RedisClient().Set(context.Background(), key, string(marshal), -1)
	//-1 means no expiration time
}

func GetPostDataViaRedis(studioID uint64) shared.PaginationData {
	studioIDStr := strconv.Itoa(int(studioID))
	key := POST_NS + studioIDStr
	var u shared.PaginationData
	bytes, _ := redis.RedisClient().Get(context.Background(), key).Bytes()
	_ = json.Unmarshal(bytes, &u)
	return u
}

func CheckIfPostStudioKeyExists(studioID uint64) bool {
	studioIDStr := strconv.Itoa(int(studioID))
	key := POST_NS + studioIDStr
	_, err := redis.RedisClient().Get(context.Background(), key).Result()
	if err != nil {
		return false
	}
	return true
}

// We'll do ot better
func BadInvalidationOfStudioPosts(studioID uint64) {
	studioIDStr := strconv.Itoa(int(studioID))
	key := POST_NS + studioIDStr
	// Delete
	redis.RedisClient().Del(context.Background(), key)
}

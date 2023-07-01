package reel

import (
	"context"
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"log"
	"strconv"
)

const REEL_STUDIO_USER_NS = "reel-users:"
const REEL_STUDIO_USER__NS = "reel-studio-users:"

// Cache Setters

func SetReelsNonStudioDataRedisWithUserID(userID uint64, data shared.NewGenericResponseV1) {
	userIDStr := strconv.Itoa(int(userID))
	key := REEL_STUDIO_USER_NS + userIDStr
	redis.RedisClient().Del(context.Background(), key)
	marshal, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
	}
	redis.RedisClient().Set(context.Background(), key, string(marshal), -1) //-1 means no expiration time
}

func (rc *reelCachingService) SetReelsWithStudioDataRedis(userID uint64, studioID uint64, data shared.NewGenericResponseV1) {
	userIDStr := strconv.Itoa(int(userID))
	studioIDStr := strconv.Itoa(int(studioID))
	key := REEL_STUDIO_USER__NS + studioIDStr + "-" + userIDStr
	redis.RedisClient().Del(context.Background(), key)
	marshal, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
	}
	redis.RedisClient().HSet(context.Background(), "cached-studio-reels:"+utils.String(studioID), userIDStr, marshal)
	//redis.RedisClient().Set(context.Background(), key, string(marshal), -1) //-1 means no expiration time
}

// Cache Data Getter

func GetReelDataViaRedisWithUserID(userID uint64) shared.NewGenericResponseV1 {
	userIDStr := strconv.Itoa(int(userID))
	key := REEL_STUDIO_USER_NS + userIDStr
	var u shared.NewGenericResponseV1
	bytes, _ := redis.RedisClient().Get(context.Background(), key).Bytes()
	_ = json.Unmarshal(bytes, &u)
	return u
}

func (rc *reelCachingService) GetReelDataWithStudioAndStudioViaRedisWithUserID(userID uint64, studioID uint64) shared.NewGenericResponseV1 {
	userIDStr := strconv.Itoa(int(userID))
	var u shared.NewGenericResponseV1
	value := rc.cache.HGet(context.Background(), "cached-studio-reels:"+utils.String(studioID), userIDStr)
	_ = json.Unmarshal([]byte(value), &u)
	return u
}

// This will check the keys Exist

func CheckIfNonStudioDataRedisWithUserIDKeyExists(userID uint64) bool {
	userIDStr := strconv.Itoa(int(userID))
	key := REEL_STUDIO_USER_NS + userIDStr
	_, err := redis.RedisClient().Get(context.Background(), key).Result()
	if err != nil {
		return false
	}
	return true
}

func (rc *reelCachingService) CheckIfReelsWithStudioDataKeyExists(userID uint64, studioID uint64) bool {
	userIDStr := strconv.Itoa(int(userID))
	//studioIDStr := strconv.Itoa(int(studioID))
	//key := REEL_STUDIO_USER__NS + studioIDStr + "-" + userIDStr
	value := rc.cache.HGet(context.Background(), "cached-studio-reels:"+utils.String(studioID), userIDStr)

	//_, err := redis.RedisClient().Get(context.Background(), key).Result()
	if value == "" {
		return false
	}
	return true
}

func (rc *reelCachingService) InvalidateReelsCachingViaStudio(studioID uint64) {
	//rc.cache.HDelete(context.Background(), "studio-reels:"+utils.String(studioID), "*")
	rc.cache.HDeleteMatching(context.Background(), "cached-studio-reels:"+utils.String(studioID), "*")
}

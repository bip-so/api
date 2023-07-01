package auth

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

func getOtpDataFromCache(userIdStr string) string {
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	val, err := rc.Get(rctx, models.RedisUserOtpNS+userIdStr).Result()
	if err == nil {
		return val
	}

	return ""
}

func deletedOtpDataFromCache(userIdStr string) {
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	rc.Del(rctx, models.RedisUserOtpNS+userIdStr)
}

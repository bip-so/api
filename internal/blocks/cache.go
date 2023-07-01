package blocks

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s blockService) MakeBranchRedisKey(idStr string) string {
	return models.REDIS_CANVAS_BRANCH_BLOCKS + idStr
}

func (s blockService) DoesKeyExists(key string) bool {
	_, err := s.cache.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// KEY does not exist
		return false
	} else if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}

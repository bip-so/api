package canvasbranch

import (
	"context"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"strconv"
)

// studioID
// branchID
// createdBlockIDs
// updatedBlockIDs
// failedToDeleteBlockIDs
//studio:branch:user:

const MainStudioLogNameSpace = "studio-blocks-log:"
const MainStudioAuditLogNameSpace = "studio-editors-log:"
const MainStudioPointsNameSpace = "studio-user-points:"

func CreateOrUpdateUserKeys(operation string, userIDStr string, studioKey string, newCount int, multiplier int) int {
	UserKey := userIDStr + "-" + operation
	userAddedKeyExists := redis.RedisClient().HGet(context.Background(), studioKey, UserKey).Val()
	if userAddedKeyExists == "" {
		nc := newCount * multiplier
		redis.RedisClient().HSet(context.Background(), studioKey, UserKey, nc)
		return 0
	} else {
		updatedAdded, _ := strconv.Atoi(userAddedKeyExists)
		updatedAdded = updatedAdded + (newCount * multiplier)
		redis.RedisClient().HSet(context.Background(), studioKey, UserKey, updatedAdded)
		return updatedAdded
	}
}

func SetStudioLogDataRedis(userID uint64, studioID uint64, branchID uint64, newBlocks int, updatedBlocks int, deletedBlocks int) {
	userIDStr := strconv.Itoa(int(userID))
	studioIDStr := strconv.Itoa(int(studioID))
	studioKey := MainStudioLogNameSpace + studioIDStr
	CreateOrUpdateUserKeys("added", userIDStr, studioKey, newBlocks, 5)
	CreateOrUpdateUserKeys("updated", userIDStr, studioKey, updatedBlocks, 1)
	// Even though we are calculating the deleted it won't be used
	CreateOrUpdateUserKeys("deleted", userIDStr, studioKey, deletedBlocks, 1)
	_ = redis.RedisClient().SAdd(context.Background(), MainStudioAuditLogNameSpace+studioIDStr, userIDStr)
}

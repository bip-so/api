package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s notificationService) GetUserIdsFromRoleIds(roleIDs *[]uint64) []uint64 {
	roles, err := App.Repo.GetRolesByID(roleIDs)
	if err != nil {
		fmt.Println(err)
	}
	userIDs := []uint64{}
	for _, role := range roles {
		for _, member := range role.Members {
			userIDs = append(userIDs, member.UserID)
		}
	}
	return userIDs
}

func (s notificationService) RemoveMentionedUserIDs(userIds, mentionUserIDs []uint64) []uint64 {
	for _, mentionId := range mentionUserIDs {
		for {
			index := utils.GetIntSliceIndex(userIds, mentionId)
			if index == nil {
				break
			}
			userIds = append(userIds[:*index], userIds[*index+1:]...)
		}
	}
	return userIds
}

func (s notificationService) GetModMentionedUserIDs(modUserIDs, mentionUserIDs []uint64) []uint64 {
	var userIDs []uint64
	for _, modId := range modUserIDs {
		for _, mentionId := range mentionUserIDs {
			if modId == mentionId {
				userIDs = append(userIDs, modId)
				break
			}
		}
	}
	return userIDs
}

func (s notificationService) GetRoughBranchNotificationsRedisKey(branchID uint64) string {
	redisKey := fmt.Sprintf("%s%d", RoughBranchNameSpace, branchID)
	return redisKey
}

func (s notificationService) AddNotificationToRedis(branchID uint64, notification *PostNotification) {
	dataStr, _ := json.Marshal(notification)
	redisKey := s.GetRoughBranchNotificationsRedisKey(branchID)
	bipredis.NewCache().HSet(context.Background(), redisKey, uuid.New().String(), dataStr)
}

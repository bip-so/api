package kafkatopics

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// CalculatePermissions and stores the data in redis.
func CalculatePermissions(msg *kafka.Message) {
	userId := utils.Uint64(string(msg.Value))

	studioPermissionsList, err := permissions.App.Service.CalculateStudioPermissions(userId)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}

	for _, studioId := range utils.Keys(studioPermissionsList) {
		collectionPermissionsList, err := permissions.App.Service.CalculateCollectionPermissions(userId, studioId)
		if err != nil {
			logger.Error(err.Error())
			KafkaConsumerError(msg, err)
			return
		}

		for _, collectionId := range utils.Keys(collectionPermissionsList) {
			canvasPermissionsList, err := permissions.App.Service.CalculateCanvasRepoPermissions(userId, studioId, collectionId)
			if err != nil {
				logger.Error(err.Error())
				KafkaConsumerError(msg, err)
				return
			}
			fmt.Println(canvasPermissionsList)

			for _, canvasId := range utils.KeysForNestedMap(canvasPermissionsList) {
				_, err := permissions.App.Service.CalculateSubCanvasRepoPermissions(userId, studioId, collectionId, canvasId)
				fmt.Println(canvasId)
				if err != nil {
					logger.Error(err.Error())
					KafkaConsumerError(msg, err)
					return
				}
			}
		}
	}
}

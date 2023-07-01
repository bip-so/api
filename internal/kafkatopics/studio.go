package kafkatopics

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func NewStudioConsumer(msg *kafka.Message) {
	// studioId := string(msg.Key)
	var stdio models.Studio
	err := json.Unmarshal(msg.Value, &stdio)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}
	isPersonal := queries.App.StudioQueries.IsPersonalStudio(stdio.ID)
	if !isPersonal {
		studio.App.StudioService.AddStudioToAlgolia(stdio.ID)
	}
	queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(stdio.CreatedByID)
}

func UpdateStudioConsumer(msg *kafka.Message) {
	// studioId := string(msg.Key)
	var stdio models.Studio
	err := json.Unmarshal(msg.Value, &stdio)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}
	isPersonal := queries.App.StudioQueries.IsPersonalStudio(stdio.ID)
	if !isPersonal {
		studio.App.StudioService.AddStudioToAlgolia(stdio.ID)
	}
	queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(stdio.CreatedByID)
	// To update the studio data in NotificationCount
	notifications.App.Service.UpdateNotificationCountAfterStudioSave(&stdio)
}

func DeletedStudioConsumer(msg *kafka.Message) {
	// studioId := string(msg.Key)
	var stdio models.Studio
	err := json.Unmarshal(msg.Value, &stdio)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}
	studio.App.StudioService.DeleteStudioFromAlgolia(stdio.ID)
	members, _ := queries.App.MemberQuery.GetMembers(map[string]interface{}{"studio_id": stdio.ID})
	for _, memberInstance := range members {
		queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(memberInstance.UserID)
	}
}

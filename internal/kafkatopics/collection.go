package kafkatopics

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

// CreateNewCollection Triggers from kafka topics When a new studio is created.
func CreateNewCollection(msg *kafka.Message) *models.Collection {
	var stdio models.Studio
	err := json.Unmarshal(msg.Value, &stdio)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return nil
	}

	collectionView := &collection.CollectionCreateValidator{
		Name:         "My new collection",
		Position:     1,
		PublicAccess: "private",
	}
	collectionInstance, err := collection.App.Controller.CreateCollectionController(collectionView, stdio.CreatedByID, stdio.ID)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return nil
	}
	return collectionInstance
}

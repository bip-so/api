package message

import (
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func CreateMessage(messages []*models.Message) error {
	return postgres.GetDB().Create(messages).Error
}

func GetMessages(userID uint64, skip int) (*[]models.Message, error) {
	var messages []models.Message
	err := postgres.GetDB().Model(&models.Message{}).Where("user_id = ? and is_used = ?", userID, false).Order("created_at desc").Preload("Author").Limit(20).Offset(skip).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return &messages, nil
}

func GetMessageByID(messageID uint64) (*models.Message, error) {
	var message models.Message
	err := postgres.GetDB().Model(&models.Message{}).Where(models.Message{BaseModel: models.BaseModel{ID: messageID}}).Preload("Author").First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func GetMessageByUUID(messageUUID string) (*models.Message, error) {
	var message models.Message
	err := postgres.GetDB().Model(&models.Message{}).Where(models.Message{BaseModel: models.BaseModel{UUID: uuid.MustParse(messageUUID)}}).Preload("Author").First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func GetMessageByRefID(refID string) (*models.Message, error) {
	var message models.Message
	err := postgres.GetDB().Model(&models.Message{}).Where(models.Message{RefID: refID}).Preload("Author").First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func MarkMessageAsUsed(messageId uint64, isUsed bool) error {
	err := postgres.GetDB().Model(&models.Message{}).Where("id = ? and is_used = ?", messageId, false).Update("is_used", isUsed).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteMessageById(userID uint64, messageId uint64) error {
	err := postgres.GetDB().Debug().Model(&models.Message{}).Where("user_id = ? and id = ?", userID, messageId).Delete(&models.Message{}).Error
	if err != nil {
		return err
	}
	return nil
}

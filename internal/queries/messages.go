package queries

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (m messageQuery) GetMessage(query map[string]interface{}) (*models.Message, error) {
	var message *models.Message
	err := postgres.GetDB().Model(models.Message{}).Where(query).First(&message).Error
	if err != nil {
		fmt.Println("error in getting message", err)
		return nil, err
	}
	return message, nil
}

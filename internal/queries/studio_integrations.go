package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q studioIntegrationQuery) GetDiscordStudioIntegration(studioId uint64) (integration models.StudioIntegration, err error) {
	condition := models.StudioIntegration{StudioID: studioId}
	condition.Type = "discord"
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).Find(&integration).Error
	return integration, err
}

func (q studioIntegrationQuery) GetStudioIntegration(query map[string]interface{}) (integration *models.StudioIntegration, err error) {
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(query).First(&integration).Error
	return
}

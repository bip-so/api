package studio_integration

import (
	"encoding/json"
	"errors"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

const (
	DISCORD_INTEGRATION_TYPE = "discord"
	SLACK_INTEGRATION_TYPE   = "slack"
)

func (r studioIntegrationsRepo) GetDiscordStudioIntegration(studioId uint64) (integration models.StudioIntegration, err error) {
	condition := models.StudioIntegration{StudioID: studioId}

	condition.Type = DISCORD_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).First(&integration).Error
	return
}

func (r studioIntegrationsRepo) GetSlackStudioIntegration(studioId uint64) (integration models.StudioIntegration, err error) {
	condition := models.StudioIntegration{StudioID: studioId}

	condition.Type = SLACK_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).First(&integration).Error
	return
}

func (r studioIntegrationsRepo) GetStudioIntegrationByDiscordTeamId(teamID string) (integration []models.StudioIntegration, err error) {
	condition := models.StudioIntegration{TeamID: teamID}

	condition.Type = DISCORD_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).Find(&integration).Error

	return
}

func (r studioIntegrationsRepo) GetActiveIntegrationForStudio(studioID uint64, integrationType string) (result *models.StudioIntegration, err error) {
	err = postgres.GetDB().Where(map[string]interface{}{
		"studio_id": studioID, // "status": true,
		"type":      integrationType,
	}).First(&result).Error
	return
}

func (r studioIntegrationsRepo) AddStudioIntegration(productId uint64, integrationType string,
	teamID string, accessToken string, extra map[string]interface{}) (*models.StudioIntegration, error) {
	extraStr, _ := json.Marshal(extra)
	var teamId string
	if integrationType == SLACK_INTEGRATION_TYPE {
		if len(teamID) != 0 {
			teamId = teamID
		} else {
			tId, ok := extra["team_id"].(string)
			if ok {
				teamId = tId
			} else {
				return nil, errors.New("team id not found")
			}
		}
	} else if integrationType == DISCORD_INTEGRATION_TYPE {
		tId, ok := extra["webhook"].(map[string]interface{})
		if !ok {
			return nil, errors.New("guild id not found")
		}
		guildId, ok := tId["guild_id"].(string)
		if ok {
			teamId = guildId
		} else {
			return nil, errors.New("guild id not found")
		}
	}
	newIntegration := models.StudioIntegration{
		Type:      integrationType,
		StudioID:  productId,
		AccessKey: accessToken,
		Status:    false,
		Extra:     extraStr,
		TeamID:    teamId,
	}
	err := postgres.GetDB().Where(models.StudioIntegration{
		Type:   integrationType,
		TeamID: teamId,
	}).Delete(&models.StudioIntegration{}).Error
	if err != nil {
		return nil, err
	}
	err = postgres.GetDB().Where(models.StudioIntegration{Type: integrationType,
		StudioID: productId, IntegrationStatus: models.StudioIntegrationPending}).FirstOrCreate(&newIntegration).Error
	return &newIntegration, err
}

func (r studioIntegrationsRepo) UpdateDiscordDmNotification(studioId uint64, status bool) error {
	err := r.db.Model(&models.Studio{}).Where("id = ?", studioId).Update("discord_notifications_enabled", status).Error
	return err
}

func (r studioIntegrationsRepo) UpdateSlackDmNotification(studioId uint64, status bool) error {
	err := r.db.Model(&models.Studio{}).Where("id = ?", studioId).Update("slack_notifications_enabled", status).Error
	return err
}

func (r studioIntegrationsRepo) DeleteStudioIntegration(studioId uint64, integrationType string) (integration models.StudioIntegration, err error) {
	condition := models.StudioIntegration{StudioID: studioId, Type: integrationType}
	err = r.db.Where(condition).Delete(&integration).Error
	return
}

func (r studioIntegrationsRepo) GetStudio(studioId uint64) (integration models.Studio, err error) {
	err = r.db.Model(&models.Studio{}).Where("id = ?", studioId).First(&integration).Error
	return
}

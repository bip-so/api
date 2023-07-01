package slack2

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func GetSlackStudioIntegration(studioId uint64) (integration []models.StudioIntegration, err error) {
	condition := models.StudioIntegration{StudioID: studioId}

	condition.Type = models.SLACK_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).Find(&integration).Error
	return
}

func (r *slackRepo) GetUserSocialAuth(providerID string) (user *models.UserSocialAuth, err error) {
	err = postgres.GetDB().Model(&models.UserSocialAuth{}).Where("provider_id = ?", providerID).Preload("User").Find(&user).Error
	return
}

func (r *slackRepo) GetStudioIntegration(teamID string) (integration *models.StudioIntegration, err error) {
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where("team_id = ?", teamID).Find(&integration).Error
	return
}

func (r *slackRepo) CreateMessage(messages *models.Message) error {
	return postgres.GetDB().Create(messages).Error
}

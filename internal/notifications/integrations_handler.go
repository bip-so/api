package notifications

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

func (s notificationService) PublishIntegrationEvent(reel *models.Reel, activity string) {
	reelString, _ := json.Marshal(reel)
	s.kafka.Publish(configs.KAFKA_INTEGRATION_EVENTS, activity, reelString)
}

func (s notificationService) SendToIntegration(reel *models.Reel) {
	integrations, _ := App.Repo.GetStudioIntegrations(reel.StudioID)
	for _, integration := range integrations {
		if integration.Type == models.DISCORD_INTEGRATION_TYPE {
			s.DiscordEventHandler(reel, &integration)
		}
	}
}

func (s notificationService) SendPostToIntegration(model *models.Post) {
	integrations, _ := App.Repo.GetStudioIntegrations(model.StudioID)
	for _, integration := range integrations {
		if integration.Type == models.DISCORD_INTEGRATION_TYPE {
			s.DiscordNewPostEventHandler(model, &integration)
		}
	}
}

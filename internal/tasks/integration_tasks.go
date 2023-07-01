package tasks

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
)

func (s taskService) SendToIntegration(ctx context.Context, task *asynq.Task) {
	var reel *models.Reel
	json.Unmarshal(task.Payload(), &reel)
	integrations, _ := notifications.App.Repo.GetStudioIntegrations(reel.StudioID)
	for _, integration := range integrations {
		if integration.Type == models.DISCORD_INTEGRATION_TYPE {
			notifications.App.Service.DiscordEventHandler(reel, &integration)
		} else if integration.Type == models.SLACK_INTEGRATION_TYPE {
			notifications.App.Service.SlackReelEventHandler(reel, &integration)
		}
	}
}

func (s taskService) SendPostToIntegration(ctx context.Context, task *asynq.Task) {
	var post *models.Post
	json.Unmarshal(task.Payload(), &post)
	integrations, _ := notifications.App.Repo.GetStudioIntegrations(post.StudioID)
	for _, integration := range integrations {
		if integration.Type == models.DISCORD_INTEGRATION_TYPE {
			notifications.App.Service.DiscordNewPostEventHandler(post, &integration)
		} else if integration.Type == models.SLACK_INTEGRATION_TYPE {
			notifications.App.Service.SlackPostEventHandler(post, &integration)
		}
	}
}
